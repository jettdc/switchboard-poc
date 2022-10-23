package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/u"
	"github.com/jettdc/switchboard/websockets"
)

func NewRoutePipeline(route config.RouteConfig, pubsubClient pubsub.PubSub) gin.HandlerFunc {
	return func(c *gin.Context) {
		listenerId := uuid.NewString()
		pipelineCtx := NewPipeContext(route, c.Params, pubsubClient, c.Request.URL.Path, listenerId)

		if err := pipelineCtx.ListenToAllTopics(); err != nil {
			u.Err(c, u.InternalServerError(err.Error()))
			return
		}

		// Upgrade request to websocket connection
		wsConnection, err := websockets.HandleConnection(c.Writer, c.Request)
		if err != nil {
			u.Err(c, u.InternalServerError("Failed to upgrade connection to websocket for route %s", c.Request.URL.Path))
			pipelineCtx.CancelFunc()
			return
		}

		// Cancel our context if there's a websocket error and continuously write incoming messages to ws
		go cancelCtxOnWSErr(pipelineCtx, wsConnection)
		go writeMessagesToWS(pipelineCtx, wsConnection)
	}
}
