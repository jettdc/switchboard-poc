package pipeline

import "github.com/jettdc/switchboard/config"

type SubscriptionTracker struct {
	SeenEndpointDescs map[*EndpointDesc]*PipeContext
}

// SubscriptionHandler is used for tracking what pipeline contexts are listening to what endpoints, primarily leveraged
// by the dynamic subscription endpoint.
type SubscriptionHandler struct {
	tracker             *SubscriptionTracker
	SubscribeRequests   chan RouteConfigWithParams
	UnsubscribeRequests chan RouteConfigWithParams
}

// CommandMessage is a type of Client->Server message used in dynamic subscriptions, where the client can request that
// certain actions ("SUBSCRIBE", "UNSUBSCRIBE") are taken on specified endpoints. A RequestId can be specified for
// associating server responses with requests.
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
