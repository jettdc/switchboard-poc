package pipeline

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/u"
	"github.com/jettdc/switchboard/websockets"
)

func NewDynamicRoutePipeline(switchboardConfig *config.Config, pubsubClient pubsub.PubSub) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Run middleware
		var err error
		for _, route := range switchboardConfig.Routes {
			for _, plugin := range route.Plugins.Middleware {
				err = (*plugin).Process(c.Request)
				if err != nil {
					u.Err(c, u.InternalServerError(err.Error()))
					return
				}
			}
		}

		// Makes sure that listen requests are idempotent through all subscribe requests
		listenerId := uuid.NewString()

		// Upgrade request to websocket connection
		wsConnection, err := websockets.HandleConnection(c.Writer, c.Request)
		if err != nil {
			u.Err(c, u.InternalServerError("Failed to upgrade connection to websocket for route %s", c.Request.URL.Path))
			return
		}

		subscriptionHandler := NewSubscriptionHandler()

		go listenForCommands(wsConnection, switchboardConfig.Routes, subscriptionHandler)
		go processCommands(subscriptionHandler, wsConnection, pubsubClient, listenerId)

		return
	}
}

func listenForCommands(wsConnection *websockets.WSConn, rcs []config.RouteConfig, sh *SubscriptionHandler) {
	for {
		var command CommandMessage
		err := wsConnection.ReadJSONSafe(&command)
		switch err {
		case nil:
			go processAction(wsConnection, command, rcs, sh)
		default:
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				wsConnection.CloseAndCancel()
				return
			} else {
				wsConnection.WriteJSONSafe(websockets.NewWSErrorMessage(fmt.Sprintf("Malformed request message: %s", err.Error()), command.RequestId))
			}
		}
	}
}

func processAction(wsConnection *websockets.WSConn, command CommandMessage, rcs []config.RouteConfig, sh *SubscriptionHandler) {
	for _, endpoint := range command.Endpoints {
		rc, err := getRouteConfigByEndpoint(endpoint.Endpoint, rcs)
		if err != nil {
			wsConnection.WriteJSONSafe(websockets.NewWSErrorMessage(fmt.Sprintf("No route configuration found for endpoint \"%s\"", endpoint.Endpoint), command.RequestId))
			continue
		} else {
			switch command.Action {
			case ActionSubscribe:
				sh.SubscribeRequests <- RouteConfigWithParams{
					rc,
					endpoint.Params,
					command.RequestId,
				}

			case ActionUnsubscribe:
				sh.UnsubscribeRequests <- RouteConfigWithParams{
					rc,
					endpoint.Params,
					command.RequestId,
				}
			default:
				u.Logger.Info(fmt.Sprintf("Command received with unknown action: %s", command.Action))
			}
		}
	}
}

func getRouteConfigByEndpoint(endpoint string, rcs []config.RouteConfig) (config.RouteConfig, error) {
	for _, rc := range rcs {
		if rc.Endpoint == endpoint {
			return rc, nil
		}
	}
	return config.RouteConfig{}, fmt.Errorf("no route config with that endpoint")
}

func processCommands(sh *SubscriptionHandler, wsConnection *websockets.WSConn, pubsubClient pubsub.PubSub, listenerId string) {
	for {
		select {
		case rc := <-sh.SubscribeRequests:
			go func() {
				ed := EndpointDesc{
					rc.RouteConfig.Endpoint,
					rc.Params,
				}

				if _, err := sh.GetPipeCtx(ed); err != nil { // Endoint desc already being handled by another pipeline
					if pipelineCtx, err := subscribeToRoute(wsConnection, rc, pubsubClient, listenerId); err != nil {
						wsConnection.WriteJSONSafe(websockets.NewWSErrorMessage(err.Error(), rc.RequestId))
					} else {
						sh.Track(pipelineCtx, &ed)
					}
				}
			}()
		case rc := <-sh.UnsubscribeRequests:
			go func() {
				if err := unsubscribeFromRoute(sh, rc); err != nil {
					wsConnection.WriteJSONSafe(websockets.NewWSErrorMessage(err.Error(), rc.RequestId))
				}
			}()
		case <-wsConnection.Ctx.Done():
			return
		}
	}
}

func subscribeToRoute(wsConnection *websockets.WSConn, route RouteConfigWithParams, pubsubClient pubsub.PubSub, listenerId string) (*PipeContext, error) {
	params := make([]gin.Param, 0)
	if route.Params != nil {
		for k, v := range *route.Params {
			params = append(params, gin.Param{Key: k, Value: v})
		}
	}

	pipelineCtx := NewPipeContextFromContext(route.RouteConfig, params, pubsubClient, route.RouteConfig.Endpoint, listenerId, wsConnection.Ctx)

	if err := pipelineCtx.ListenToAllTopics(); err != nil {
		return nil, err
	}

	// Continuously write incoming messages to ws
	go writeMessagesToWS(pipelineCtx, wsConnection)
	return pipelineCtx, nil
}

func unsubscribeFromRoute(sh *SubscriptionHandler, route RouteConfigWithParams) error {
	pipelineCtx, err := sh.GetPipeCtx(EndpointDesc{
		route.RouteConfig.Endpoint,
		route.Params,
	})
	if err != nil {
		return err
	}
	sh.CancelAndDeleteEntry(pipelineCtx)
	return nil
}
