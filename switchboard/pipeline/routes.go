package pipeline

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/websockets"
)

func NewRoutePipeline(route config.RouteConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Route messages through any middleware if it exists

		// TODO: Should validate topic before, when loading in config?
		parameterizedTopic, err := ParameterizeTopic(route.Topic, c.Params)
		if err != nil {
			c.Writer.WriteHeader(500)
			c.Writer.WriteString("Invalid pubsub path.")
			c.Done()
			return
		}

		ctx, cancelFunc := context.WithCancel(context.Background())
		messages, err := pubsub.Redis.Subscribe(ctx, parameterizedTopic)
		if err != nil {
			cancelFunc()
			c.Writer.WriteHeader(500)
			c.Writer.WriteString("Failed to subscribe to redis topic.")
			c.Done()
		}

		// Upgrade request to websocket connection
		wsConnection, err := websockets.HandleConnection(c.Writer, c.Request)
		if err != nil {
			c.Writer.WriteHeader(500)
			c.Writer.WriteString("Failed to upgrade connection to websocket.")
			c.Done()
			cancelFunc()
			return
		}

		// Cancel our pubsub subscription if there's a websocket error
		go func() {
			for {
				_, _, err = wsConnection.ReadMessage()
				if err != nil {
					cancelFunc()
					return
				}
			}
		}()

		// Continuously write messages to socket
		go func() {
			for {
				select {
				case msg := <-messages:
					// TODO route incoming messages through enrichment plugins, if any
					wsConnection.WriteJSON(msg)
				case <-ctx.Done():
					return
				}
			}
		}()

		return
	}
}
