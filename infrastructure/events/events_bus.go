package events

import (
	"context"
	"fmt"
)

type HandlerFnc = func(ctx context.Context, message interface{})

type EventBus struct {
	Bus      chan interface{}
	handlers map[string]HandlerFnc
}

// Serve implements IEventBus.
func (e EventBus) Serve(ctx context.Context) {
	for {
		select {
		case evnet := <-e.Bus:
			{
				fmt.Printf("EVENT: %+v\n", evnet)
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

	e.handlers[topic] = handler

}

type IEventBus interface {
	Dispatch(topic string, message interface{})
	Listen(topic string, handler HandlerFnc)
	Serve(ctx context.Context)
}

func NewEventBus() IEventBus {
	return EventBus{
		Bus:      make(chan interface{}),
		handlers: map[string]func(ctx context.Context, message interface{}){},
	}
}
