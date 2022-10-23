package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jettdc/switchboard/u"
	"strings"
)

type PubSub interface {
	Connect() error
	Subscribe(ctx context.Context, topic string, listenerId string) (chan ForwardedMessage, error)
}

type ForwardedMessage struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

func (p ForwardedMessage) String() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("invalid pubsub message json")
	}
	return string(b), nil
}

func GetPubSubClient(provider string) (PubSub, error) {
	switch strings.ToLower(provider) {
	case "redis":
		u.Logger.Info("Switchboard configured to use Redis as pubsub provider.")
		return Redis, nil
	default:
		return nil, fmt.Errorf("cannot find supported pubsub provider with name \"%s\"", provider)
	}
}
