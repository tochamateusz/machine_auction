package events

import (
	"context"
	"github.com/rs/zerolog/log"
)

type HandlerFnc = func(ctx context.Context, message interface{})

type Bus struct {
	Topic   string
	Message interface{}
}

type EventBus struct {
	Bus      chan Bus
	handlers map[string][]HandlerFnc
}

// Serve implements IEventBus.
func (e EventBus) Serve(ctx context.Context) {
	for {
		select {
		case event := <-e.Bus:
			{
				switch event.Topic {
				case "auction.founded":
					{
						handlers, ok := e.handlers[event.Topic]
						if ok == false {
							log.Warn().Msgf("handler for event: %s not register", event.Topic)
							continue
						}
						for _, v := range handlers {
							go v(ctx, event.Message)
						}
					}

				case "auctions.founded":
					{
						handlers, ok := e.handlers[event.Topic]
						if ok == false {
							log.Warn().Msgf("handler for event: %s not register", event.Topic)
							continue
						}
						for _, v := range handlers {
							go v(ctx, event.Message)
						}
					}

				default:
					{
						continue
					}
				}

			}
		}
	}
}

// Dispatch implements IEventBus.
func (e EventBus) Dispatch(topic string, message interface{}) {
	e.Bus <- struct {
		Topic   string
		Message interface{}
	}{
		Topic:   topic,
		Message: message,
	}
}

// Listen implements IEventBus.
func (e EventBus) Listen(topic string, handler HandlerFnc) {
	e.handlers[topic] = append(e.handlers[topic], handler)

}

type IEventBus interface {
	Dispatch(topic string, message interface{})
	Listen(topic string, handler HandlerFnc)
	Serve(ctx context.Context)
}

func NewEventBus() IEventBus {
	return EventBus{
		Bus:      make(chan Bus, 3),
		handlers: map[string][]func(ctx context.Context, message interface{}){},
	}
}
