package transport

import "context"

// Message payload, including the sender ClientID.
type Message struct {
	ClientID string
	Name     string
	Data     any
}

// MessageHandlerFn callback
type MessageHandlerFn func(*Message)

// UnsubscribeFn to be invoked when a subscription is no longer needed
type UnsubscribeFn func()

// Transporter provides the ability to send and receive messages to a set of clients.
type Transporter interface {
	ClientIDs(context.Context) ([]string, error)
	Publish(context.Context, string, any) error
	Subscribe(context.Context, string, MessageHandlerFn) (UnsubscribeFn, error)
}
