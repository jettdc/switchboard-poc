package pipeline

import (
	"fmt"
	"github.com/jettdc/switchboard/u"
	"github.com/jettdc/switchboard/websockets"
)

func cancelCtxOnWSErr(pipelineContext *PipeContext, wsConn *websockets.WSConn) {
	for {
		_, _, err := wsConn.ReadMessageSafe()
		if err != nil {
			pipelineContext.CancelFunc()
			return
		}
	}
}

func writeMessagesToWS(pipeContext *PipeContext, wsConn *websockets.WSConn) {
	for {
		select {
		case msg := <-pipeContext.AllMessages:
			// TODO route incoming messages through enrichment plugins, if any
			j, err := msg.String()
			if err != nil {
				fmt.Printf("Error converting pubsub message to json.")
			}
			wsConn.WriteJSONSafe(j)
		case <-wsConn.Ctx.Done():
			u.Logger.Info(fmt.Sprintf("Client disconnected from websocket at %s.", pipeContext.ResolvedEndpoint))
			return
		}
	}
}
