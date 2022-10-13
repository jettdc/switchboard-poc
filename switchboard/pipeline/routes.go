package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/websockets"
)

func NewRoutePipeline(route config.RouteConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Route messages through any middleware if it exists

		// TODO: subscribe to topic

		// Upgrade request to websocket connection
		wsConnection, err := websockets.HandleConnection(c.Writer, c.Request)
		if err != nil {
			c.Writer.WriteHeader(500)
			c.Writer.WriteString("Failed to upgrade connection to websocket.")
			c.Done()
			return
		}

		// TODO route incoming messages through enrichment plugins, if any

		// TODO Forward messages from redis to websocket

		var res = make(map[string]string)
		res["test"] = "test2"
		wsConnection.WriteJSON(res)
		return
	}
}
