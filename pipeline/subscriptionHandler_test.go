package pipeline

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/mock"
	"testing"
)

func TestSubscriptionTracker_TrackEndpointDesc_DoesntExist(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	err := subscriptionHandler.Track(pc, &ed)
	if err != nil {
		t.Errorf("shouldn't return an error")
	}
	if subscriptionHandler.tracker.SeenEndpointDescs[&ed] != pc {
		t.Errorf("tracking an endpoint description for the first time should track the current pipeline context")
	}

	cf()
}

func TestSubscriptionTracker_TrackEndpointDesc_Exists(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	err := subscriptionHandler.Track(pc, &ed)
	err = subscriptionHandler.Track(pc, &ed2)
	if err == nil {
		t.Errorf("should return an error")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_Nonetracked(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	_, err := subscriptionHandler.GetPipeCtx(ed)

	if err == nil {
		t.Errorf("should return an error if there are no pipelines associated with the requested endpoint desc")
	}
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_TrackedNoParams(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionHandler.tracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	pipelineCtx, err := subscriptionHandler.GetPipeCtx(ed2)

	if err != nil {
		t.Errorf("should not return an error")
	}

	if pipelineCtx != pc {
		t.Errorf("should return the existing pipeline if a same-behaving endpoint is already being tracked")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_TrackedWithParams(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	p1 := make(map[string]string)
	p1["param"] = "value"
	ed := EndpointDesc{Endpoint: "/firsttest", Params: &p1}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionHandler.tracker.SeenEndpointDescs[&ed] = pc

	p2 := make(map[string]string)
	p2["param"] = "value"
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p2}

	pipelineCtx, err := subscriptionHandler.GetPipeCtx(ed2)

	if err != nil {
		t.Errorf("should not return an error")
	}

	if pipelineCtx != pc {
		t.Errorf("should return the existing pipeline if a same-behaving endpoint is already being tracked")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UntrackedNoParams(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionHandler.tracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest2", Params: nil}

	pipelineCtx, err := subscriptionHandler.GetPipeCtx(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("endpoint in second desc doesn't exist, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UnrackedOneParams(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	p1 := make(map[string]string)
	p1["param"] = "value"
	ed := EndpointDesc{Endpoint: "/firsttest", Params: &p1}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionHandler.tracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	pipelineCtx, err := subscriptionHandler.GetPipeCtx(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UnrackedOneParamsOther(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	p1 := make(map[string]string)
	p1["param"] = "value"
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionHandler.tracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p1}

	pipelineCtx, err := subscriptionHandler.GetPipeCtx(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UntrackedBothParams(t *testing.T) {
	subscriptionHanlder := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	p1 := make(map[string]string)
	p1["param"] = "value"
	ed := EndpointDesc{Endpoint: "/firsttest", Params: &p1}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionHanlder.tracker.SeenEndpointDescs[&ed] = pc

	p2 := make(map[string]string)
	p2["param"] = "value2"
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p2}

	pipelineCtx, err := subscriptionHanlder.GetPipeCtx(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UntrackedBothParamsDiffKeys(t *testing.T) {
	subscriptionHandler := NewSubscriptionHandler()

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	p1 := make(map[string]string)
	p1["param"] = "value"
	ed := EndpointDesc{Endpoint: "/firsttest", Params: &p1}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionHandler.tracker.SeenEndpointDescs[&ed] = pc

	p2 := make(map[string]string)
	p2["param2"] = "value"
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p2}

	pipelineCtx, err := subscriptionHandler.GetPipeCtx(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}
