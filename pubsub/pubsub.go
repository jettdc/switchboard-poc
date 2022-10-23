// For encapsulating the pubsub provider logic to allow easy switching
package pubsub

import (
	"context"
	"fmt"
	"github.com/jettdc/switchboard/pubsub/listen_groups"
	"github.com/jettdc/switchboard/u"
	"strings"
)

// PubSub is the generic interface for interacting with various pubsub providers, such as redis.
type PubSub interface {
	Connect() error
	Subscribe(ctx context.Context, topic string, listenerId string) (chan listen_groups.ForwardedMessage, error)
}

// GetPubSubClient matches the requested provider with the available providers, and returns an instance of it.
// Providers must be explicitly entered into the source code
func GetPubSubClient(provider string) (PubSub, error) {
	switch strings.ToLower(provider) {
	case "redis":
		u.Logger.Info("Switchboard configured to use Redis as pubsub provider.")
		return Redis, nil
	default:
		return nil, fmt.Errorf("cannot find supported pubsub provider with name \"%s\"", provider)
	}
}
