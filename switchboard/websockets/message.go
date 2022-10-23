package websockets

import (
	"encoding/json"
	"fmt"
	"github.com/jettdc/switchboard/u"
)

const (
	ForwardedMessage string = "FORWARDED_MESSAGE"
	Response                = "RESPONSE"
	Error                   = "ERROR"
)

type Message struct {
	Endpoint    string      `json:"endpoint"`
	MessageType string      `json:"messageType"`
	Message     interface{} `json:"message"`
	RequestId   *string     `json:"requestId,omitempty"`
}

type ErrorPayload struct {
	Error string `json:"error"`
}

func NewWSErrorMessage(msg string, requestId *string) Message {
	u.Logger.Debug(fmt.Sprintf("Writing err:  %s", msg))
	return Message{
		"/multi",
		Error,
		ErrorPayload{msg},
		requestId,
	}
}

func (p Message) String() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("invalid message json")
	}
	return string(b), nil
}
