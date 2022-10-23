package listen_groups

import (
	"encoding/json"
	"fmt"
)

// ForwardedMessage creates a standardized message type for all pubsub providers to pass along.
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
