package pipeline

import "github.com/jettdc/switchboard/config"

type SubscriptionTracker struct {
	SeenEndpointDescs map[*EndpointDesc]*PipeContext
}

type SubscriptionHandler struct {
	tracker             *SubscriptionTracker
	SubscribeRequests   chan RouteConfigWithParams
	UnsubscribeRequests chan RouteConfigWithParams
}

type CommandMessage struct {
	Action    string         `json:"action"`
	Endpoints []EndpointDesc `json:"endpoints"`
	RequestId *string        `json:"requestId,omitempty"`
}

type EndpointDesc struct {
	Endpoint string             `json:"endpoint"`
	Params   *map[string]string `json:"params,omitempty"`
}

type RouteConfigWithParams struct {
	RouteConfig config.RouteConfig
	Params      *map[string]string
	RequestId   *string
}

const (
	ActionSubscribe   string = "SUBSCRIBE"
	ActionUnsubscribe        = "UNSUBSCRIBE"
)
