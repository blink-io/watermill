package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// BookRoomHandler is a command handler, which handles BookRoom command and emits RoomBooked.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be retuerned.
type BookRoomHandler struct {
	eventBus *cqrs.EventBus
}

func (b BookRoomHandler) HandlerName() string {
	return "BookRoomHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b BookRoomHandler) NewCommand() interface{} {
	return &BookRoom{}
}

func (b BookRoomHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*BookRoom)

	// some random price, in production you probably will calculate in wiser way
	price := (rand.Int63n(40) + 1) * 10

	log.Printf(
		"Booked %s for %s from %s to %s",
		cmd.RoomId,
		cmd.GuestName,
		time.Unix(cmd.StartDate.Seconds, int64(cmd.StartDate.Nanos)),
		time.Unix(cmd.EndDate.Seconds, int64(cmd.EndDate.Nanos)),
	)

	// RoomBooked will be handled by OrderBeerOnRoomBooked event handler,
	// in future RoomBooked may be handled by multiple event handler
	if err := b.eventBus.Publish(ctx, &RoomBooked{
		ReservationId: watermill.NewUUID(),
		RoomId:        cmd.RoomId,
		GuestName:     cmd.GuestName,
		Price:         price,
		StartDate:     cmd.StartDate,
		EndDate:       cmd.EndDate,
	}); err != nil {
		return err
	}

	return nil
}

// OrderBeerOnRoomBooked is a event handler, which handles RoomBooked event and emits OrderBeer command.
type OrderBeerOnRoomBooked struct {
	commandBus *cqrs.CommandBus
}

func (o OrderBeerOnRoomBooked) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "OrderBeerOnRoomBooked"
}

func (OrderBeerOnRoomBooked) NewEvent() interface{} {
	return &RoomBooked{}
}

func (o OrderBeerOnRoomBooked) Handle(ctx context.Context, e interface{}) error {
	event := e.(*RoomBooked)

	orderBeerCmd := &OrderBeer{
		RoomId: event.RoomId,
		Count:  rand.Int63n(10) + 1,
	}

	return o.commandBus.Send(ctx, orderBeerCmd)
}

// OrderBeerHandler is a command handler, which handles OrderBeer command and emits BeerOrdered.
// BeerOrdered is not handled by any event handler, but we may use persistent Pub/Sub to handle it in the future.
type OrderBeerHandler struct {
	eventBus *cqrs.EventBus
}

func (o OrderBeerHandler) HandlerName() string {
	return "OrderBeerHandler"
}

func (o OrderBeerHandler) NewCommand() interface{} {
	return &OrderBeer{}
}

func (o OrderBeerHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*OrderBeer)

	if rand.Int63n(10) == 0 {
		// sometimes there is no beer left, command will be retried
		return errors.Errorf("no beer left for room %s, please try later", cmd.RoomId)
	}

	if err := o.eventBus.Publish(ctx, &BeerOrdered{
		RoomId: cmd.RoomId,
		Count:  cmd.Count,
	}); err != nil {
		return err
	}

	log.Printf("%d beers ordered to room %s", cmd.Count, cmd.RoomId)
	return nil
}

// BookingsFinancialReport is a read model, which calculates how much money we may earn from bookings.
// Like OrderBeerOnRoomBooked, it listens for RoomBooked event.
//
// This implementation is just writing to the memory. In production, you will probably will use some persistent storage.
type BookingsFinancialReport struct {
	handledBookings map[string]struct{}
	totalCharge     int64
	lock            sync.Mutex
}

func NewBookingsFinancialReport() *BookingsFinancialReport {
	return &BookingsFinancialReport{handledBookings: map[string]struct{}{}}
}

func (b BookingsFinancialReport) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "BookingsFinancialReport"
}

func (BookingsFinancialReport) NewEvent() interface{} {
	return &RoomBooked{}
}

func (b *BookingsFinancialReport) Handle(ctx context.Context, e interface{}) error {
	// Handle may be called concurrently, so it need to be thread safe.
	b.lock.Lock()
	defer b.lock.Unlock()

	event := e.(*RoomBooked)

	// When we are using Pub/Sub which doesn't provide exactly-once delivery semantics, we need to deduplicate messages.
	// GoChannel Pub/Sub provides exactly-once delivery,
	// but let's make this example ready for other Pub/Sub implementations.
	if _, ok := b.handledBookings[event.ReservationId]; ok {
		return nil
	}
	b.handledBookings[event.ReservationId] = struct{}{}

	b.totalCharge += event.Price

	fmt.Printf(">>> Already booked rooms for $%d\n", b.totalCharge)
	return nil
}

var amqpAddress = "amqp://guest:guest@rabbitmq:5672/"

func main() {
	logger := watermill.NewStdLogger(false, false)
	cqrsMarshaler := cqrs.ProtobufMarshaler{}

	// You can use any Pub/Sub implementation from here: https://watermill.io/pubsubs/
	// Detailed RabbitMQ implementation: https://watermill.io/pubsubs/amqp/
	// Commands will be send to queue, because they need to be consumed once.
	commandsAMQPConfig := amqp.NewDurableQueueConfig(amqpAddress)
	commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}
	commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}

	// Events will be published to PubSub configured Rabbit, because they may be consumed by multiple consumers.
	// (in that case BookingsFinancialReport and OrderBeerOnRoomBooked).
	eventsPublisher, err := amqp.NewPublisher(amqp.NewDurablePubSubConfig(amqpAddress, nil), logger)
	if err != nil {
		panic(err)
	}

	// CQRS is built on messages router. Detailed documentation: https://watermill.io/docs/messages-router/
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	// Simple middleware which will recover panics from event or command handlers.
	// More about router middlewares you can find in the documentation:
	// https://watermill.io/docs/messages-router/#middleware
	//
	// List of available middlewares you can find in message/router/middleware.
	router.AddMiddleware(middleware.Recoverer)

	commandsConfig := cqrs.CommandConfig{
		GeneratePublishTopic: func(params cqrs.GenerateCommandPublishTopicParams) (string, error) {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return params.CommandName, nil
		},
		GenerateHandlerSubscribeTopic: func(params cqrs.GenerateCommandHandlerSubscribeTopicParams) (string, error) {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return params.CommandName, nil
		},
		SubscriberConstructor: func(params cqrs.CommandsSubscriberConstructorParams) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return commandsSubscriber, nil
		},
		OnSend: func(params cqrs.OnCommandSendParams) error {
			logger.Info("Sending command", watermill.LogFields{
				"command_name": params.CommandName,
			})

			params.Message.Metadata.Set("sent_at", time.Now().String())

			return nil
		},
		OnHandle: func(params cqrs.OnCommandHandleParams) error {
			start := time.Now()

			err := params.Handler.Handle(params.Message.Context(), params.Command)

			logger.Info("Command handled", watermill.LogFields{
				"command_name": params.CommandName,
				"duration":     time.Since(start),
				"err":          err,
			})

			return err
		},
		Marshaler: cqrsMarshaler,
		Logger:    logger,
	}

	commandBus, err := cqrs.NewCommandBusWithConfig(commandsPublisher, commandsConfig)
	if err != nil {
		panic(err)
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(commandsConfig)
	if err != nil {
		panic(err)
	}

	eventsConfig := cqrs.EventConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "events", nil

			// we can also use topic per event type
			// return params.EventName, nil
		},
		GenerateHandlerGroupSubscribeTopic: func(params cqrs.GenerateEventHandlerGroupTopicParams) (string, error) {
			return "events", nil
		},
		GroupSubscriberConstructor: func(params cqrs.EventsGroupSubscriberConstructorParams) (message.Subscriber, error) {
			config := amqp.NewDurablePubSubConfig(
				amqpAddress,
				amqp.GenerateQueueNameTopicNameWithSuffix(params.EventGroupName),
			)

			return amqp.NewSubscriber(config, logger)
		},
		OnPublish: func(params cqrs.OnEventSendParams) error {
			logger.Info("Publishing event", watermill.LogFields{
				"event_name": params.EventName,
			})

			params.Message.Metadata.Set("published_at", time.Now().String())

			return nil
		},
		OnHandle: func(params cqrs.OnEventHandleParams) error {
			start := time.Now()

			err := params.Handler.Handle(params.Message.Context(), params.Event)

			logger.Info("Event handled", watermill.LogFields{
				"event_name": params.EventName,
				"duration":   time.Since(start),
				"err":        err,
			})

			return err
		},
		Marshaler: cqrsMarshaler,
		Logger:    logger,
	}

	eventBus, err := cqrs.NewEventBusWithConfig(eventsPublisher, eventsConfig)
	if err != nil {
		panic(err)
	}

	commandProcessor.AddHandler(
		BookRoomHandler{eventBus},
		OrderBeerHandler{eventBus},
	)

	err = commandProcessor.AddHandlersToRouter(router)
	if err != nil {
		panic(err)
	}

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(eventsConfig)
	if err != nil {
		panic(err)
	}

	err = eventProcessor.AddHandlersGroup(
		"events",
		OrderBeerOnRoomBooked{commandBus},

		NewBookingsFinancialReport(),

		cqrs.NewGroupEventHandler(func(ctx context.Context, event *BeerOrdered) error {
			logger.Info("Beer ordered", watermill.LogFields{
				"room_id": event.RoomId,
			})
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	if err := eventProcessor.AddHandlersToRouter(router); err != nil {
		panic(err)
	}

	// publish BookRoom commands every second to simulate incoming traffic
	go publishCommands(commandBus)

	// processors are based on router, so they will work when router will start
	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}

func publishCommands(commandBus *cqrs.CommandBus) func() {
	i := 0
	for {
		i++

		startDate, err := ptypes.TimestampProto(time.Now())
		if err != nil {
			panic(err)
		}

		endDate, err := ptypes.TimestampProto(time.Now().Add(time.Hour * 24 * 3))
		if err != nil {
			panic(err)
		}

		bookRoomCmd := &BookRoom{
			RoomId:    fmt.Sprintf("%d", i),
			GuestName: "John",
			StartDate: startDate,
			EndDate:   endDate,
		}
		if err := commandBus.Send(context.Background(), bookRoomCmd); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}
