package pipeline

import (
	"fmt"
	"github.com/jettdc/switchboard/u"
	"github.com/jettdc/switchboard/websockets"
)

func writeMessagesToWS(pipeContext *PipeContext, wsConn *websockets.WSConn) {
	for {
		select {
		case msg := <-pipeContext.AllMessages:
			wsConn.WriteJSONSafe(msg)
		case <-wsConn.Ctx.Done():
			u.Logger.Info(fmt.Sprintf("Client disconnected from websocket at %s.", pipeContext.ResolvedEndpoint))
			return
		}
	}
}
