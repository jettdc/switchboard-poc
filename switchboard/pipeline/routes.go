package pipeline

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/u"
	"github.com/jettdc/switchboard/websockets"
)

func NewRoutePipeline(route config.RouteConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Route connection through any middleware if it exists

		ctx, cancelFunc := context.WithCancel(context.Background())
		allMessages := make(chan pubsub.Message, 8)

		// Subscribe to all specified topics
		// route all messages for this route into a single channel
		for _, topic := range route.Topics {
			if err := listenOnTopic(topic, c.Params, allMessages, ctx); err != nil {
				u.Err(c, u.InternalServerError("could not subscribe to topic %s", topic))
				cancelFunc()
				return
			}
		}

		// Upgrade request to websocket connection
		wsConnection, err := websockets.HandleConnection(c.Writer, c.Request)
		if err != nil {
			u.Err(c, u.InternalServerError("Failed to upgrade connection to websocket for route %s", c.Request.URL.Path))
			cancelFunc()
			return
		}

		// Cancel our context if there's a websocket error and continuously write incoming messages to ws
		go cancelCtxOnWSErr(wsConnection, cancelFunc)
		go writeMessagesToWS(allMessages, wsConnection, ctx)

		return
	}
}

// Subscribe to a topic and forward all messages to the single channel
func listenOnTopic(topic string, params gin.Params, allMessages chan pubsub.Message, ctx context.Context) error {
	// /example/topic/:id -> /example/topic/3
	// Don't need to check for error, topics are validated on config load
	parameterizedTopic, _ := config.ParameterizeTopic(topic, params)

	topicMessages, err := pubsub.Redis.Subscribe(ctx, parameterizedTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-topicMessages:
				allMessages <- msg
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func cancelCtxOnWSErr(wsConnection *websocket.Conn, cancelFunc context.CancelFunc) {
	for {
		_, _, err := wsConnection.ReadMessage()
		if err != nil {
			cancelFunc()
			return
		}
	}
}

func writeMessagesToWS(messages chan pubsub.Message, wsConnection *websocket.Conn, ctx context.Context) {
	for {
		select {
		case msg := <-messages:
			// TODO route incoming messages through enrichment plugins, if any
			j, err := msg.String()
			if err != nil {
				fmt.Printf("Error converting pubsub message to json.")
			}

			wsConnection.WriteJSON(j)
		case <-ctx.Done():
			wsConnection.Close()
			return
		}
	}
}
