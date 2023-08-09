package ably

import (
	"context"
	"fmt"

	"github.com/KyleBanks/ably-test/server/drivers/transport"
	"github.com/ably/ably-go/ably"
)

type client struct {
	channel *ably.RealtimeChannel
}

// NewTransport initialises an Ably transporter implementation, allowing
// for publishing and receiving messages over an Ably channel.
func NewTransport(
	ctx context.Context,
	apiKey, channelName string,
) (transport.Transporter, error) {
	conn, err := ably.NewRealtime(ably.WithKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to initialise Ably transport: %w", err)
	}

	c := &client{
		channel: conn.Channels.Get(channelName),
	}

	conn.Connect()
	return c, nil
}

func (c *client) ClientIDs(ctx context.Context) ([]string, error) {
	clients, err := c.channel.Presence.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get presence: %w", err)
	}

	var clientIDs []string
	for _, client := range clients {
		clientIDs = append(clientIDs, client.ClientID)
	}

	return clientIDs, nil
}

func (c *client) Publish(ctx context.Context, name string, data any) error {
	return c.channel.Publish(ctx, name, data)
}

func (c *client) Subscribe(ctx context.Context, name string, fn transport.MessageHandlerFn) (transport.UnsubscribeFn, error) {
	return c.channel.Subscribe(ctx, name, func(m *ably.Message) {
		fn(&transport.Message{
			ClientID: m.ClientID,
			Name:     name,
			Data:     m.Data,
		})
	})
}
