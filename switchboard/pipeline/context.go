package pipeline

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/websockets"
)

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

func NewPipeContext(route config.RouteConfig, params gin.Params, pubsubClient pubsub.PubSub, path string, listenerId string) *PipeContext {
	newCtx, cancelFunc := context.WithCancel(context.Background())
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

// Subscribe to all specified topics
// route all messages for this route into a single channel
func (p *PipeContext) ListenToAllTopics() error {
	for _, topic := range p.RouteConfig.Topics {
		if err := p.listenOnTopic(topic); err != nil {
			p.CancelFunc()
			return fmt.Errorf("could not subscribe to topic %s", topic)
		}
	}
	return nil
}

// Subscribe to a topic and forward all messages to the single channel
func (p *PipeContext) listenOnTopic(topic string) error {
	// TODO: Make sure not listening twice

	// /example/topic/:id -> /example/topic/3
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
