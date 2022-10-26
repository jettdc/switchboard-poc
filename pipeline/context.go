// Package pipeline contains the logic for user subscriptions.
//
// A "pipeline" consists of:
//   - A user connecting on a websocket
//   - The server subscribing to a topic, forwarding those messages
//   - Plugins being utilized
package pipeline

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/websockets"
)

// PipeContext serves as an abstraction for a connection to a pubsub provider.
// It represents a single listener, an important concept for [pubsub/listen_groups]. It provides a convenient
// centralization and abstraction for subscribing to topics, forwarding those messages to a channel, and handling the
// lifecycle of the listener.
type PipeContext struct {
	Ctx              context.Context
	CancelFunc       context.CancelFunc
	RouteConfig      config.RouteConfig
	ResolvedEndpoint string
	RouteParams      gin.Params
	AllMessages      chan websockets.Message
	PubSub           pubsub.PubSub
	ListenerId       string
}

func NewPipeContextFromContext(route config.RouteConfig, params gin.Params, pubsubClient pubsub.PubSub, path string, listenerId string, ctx context.Context) *PipeContext {
	newCtx, cancelFunc := context.WithCancel(ctx)
	pipelineCtx := PipeContext{
		ResolvedEndpoint: path,
		RouteConfig:      route,
		RouteParams:      params,
		AllMessages:      make(chan websockets.Message, 8),
		Ctx:              newCtx,
		CancelFunc:       cancelFunc,
		PubSub:           pubsubClient,
		ListenerId:       listenerId,
	}
	return &pipelineCtx
}

// ListenToAllTopics subscribes to all topics listed in the [PipeContext] RouteConfig, as well as parameterizes them.
// All messages are then forwarded to the [PipeContext] message channel.
func (p *PipeContext) ListenToAllTopics() error {
	for _, topic := range p.RouteConfig.Topics {
		if err := p.listenOnTopic(topic); err != nil {
			p.CancelFunc()
			return fmt.Errorf("could not subscribe to topic %s", topic)
		}
	}
	return nil
}

func (p *PipeContext) listenOnTopic(topic string) error {

	// /examples/topic/:id -> /examples/topic/3
	// Don't need to check for error, topics are validated on config load
	parameterizedTopic, _ := config.ParameterizeTopic(topic, p.RouteParams)

	topicMessages, err := p.PubSub.Subscribe(p.Ctx, parameterizedTopic, p.ListenerId)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-topicMessages:
				toForward := websockets.Message{
					p.ResolvedEndpoint,
					websockets.ForwardedMessage,
					msg,
					nil,
				}
				p.AllMessages <- toForward
			case <-p.Ctx.Done():
				return
			}
		}
	}()

	return nil
}
