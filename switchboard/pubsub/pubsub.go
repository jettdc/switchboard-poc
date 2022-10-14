package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jettdc/switchboard/common"
)

type PubSub interface {
	Connect() error
	Subscribe(ctx context.Context, topic string) (chan Message, error)
}

type Message struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

func (p Message) String() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("invalid pubsub message json")
	}
	return string(b), nil
}

func GetPubSubClient() (PubSub, error) {
	provider := common.GetEnvWithDefault("PUBSUB_PROVIDER", "redis")
	switch provider {
	case "redis":
		return Redis, nil
	default:
		return nil, fmt.Errorf("cannot find pubsub provider with name %s", provider)
	}
}
