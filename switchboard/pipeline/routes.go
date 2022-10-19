package pipeline

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/websockets"
)

func NewRoutePipeline(route config.RouteConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Route messages through any middleware if it exists

		// TODO: Should validate topic before, when loading in config?

		ctx, cancelFunc := context.WithCancel(c.Request.Context())
		allMessages := make(chan pubsub.Message, 8)

		// route all target topic messages into a single channel
		for _, topic := range route.Topics {
			// Don't need to check for error, topics are validated on config load
			parameterizedTopic, _ := config.ParameterizeTopic(topic, c.Params)

			topicMessages, err := pubsub.Redis.Subscribe(ctx, parameterizedTopic)
			if err != nil {
				cancelFunc()
				c.Writer.WriteHeader(500)
				c.Writer.WriteString("Failed to subscribe to redis topic.")
				c.Done()
			}

			// demux messages
			go func() {
				for msg := range topicMessages {
					allMessages <- msg
				}
			}()
		}

		// Upgrade request to websocket connection
		wsConnection, err := websockets.HandleConnection(c.Writer, c.Request)
		if err != nil {
			c.Writer.WriteHeader(500)
			c.Writer.WriteString("Failed to upgrade connection to websocket.")
			c.Done()
			cancelFunc()
			close(allMessages)
			return
		}

		// Cancel our pubsub subscription if there's a websocket error
		go func() {
			for {
				_, _, err = wsConnection.ReadMessage()
				if err != nil {
					cancelFunc()
					close(allMessages)
					return
				}
			}
		}()

		// Continuously write messages to socket
		go func() {
			for {
				select {
				case msg := <-allMessages:
					// TODO route incoming messages through enrichment plugins, if any
					j, err := msg.String()
					if err != nil {
						fmt.Printf("Error converting pubsub message to json.")
					}

					wsConnection.WriteJSON(j)
				case <-ctx.Done():
					close(allMessages)
					return
				}
			}
		}()

		return
	}
}
