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
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	err := subscriptionTracker.TrackEndpointDesc(pc, &ed)
	if err != nil {
		t.Errorf("shouldn't return an error")
	}
	if subscriptionTracker.SeenEndpointDescs[&ed] != pc {
		t.Errorf("tracking an endpoint description for the first time should track the current pipeline context")
	}

	cf()
}

func TestSubscriptionTracker_TrackEndpointDesc_Exists(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	err := subscriptionTracker.TrackEndpointDesc(pc, &ed)
	err = subscriptionTracker.TrackEndpointDesc(pc, &ed2)
	if err == nil {
		t.Errorf("should return an error")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_Nonetracked(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	_, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed)

	if err == nil {
		t.Errorf("should return an error if there are no pipelines associated with the requested endpoint desc")
	}
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_TrackedNoParams(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionTracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	pipelineCtx, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed2)

	if err != nil {
		t.Errorf("should not return an error")
	}

	if pipelineCtx != pc {
		t.Errorf("should return the existing pipeline if a same-behaving endpoint is already being tracked")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_TrackedWithParams(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

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
	subscriptionTracker.SeenEndpointDescs[&ed] = pc

	p2 := make(map[string]string)
	p2["param"] = "value"
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p2}

	pipelineCtx, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed2)

	if err != nil {
		t.Errorf("should not return an error")
	}

	if pipelineCtx != pc {
		t.Errorf("should return the existing pipeline if a same-behaving endpoint is already being tracked")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UntrackedNoParams(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mock.NewMockPubSub(ctrl)
	c, cf := context.WithCancel(context.Background())
	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	// One with an endpoint that behaves the same is already being tracked
	subscriptionTracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest2", Params: nil}

	pipelineCtx, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("endpoint in second desc doesn't exist, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UnrackedOneParams(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

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
	subscriptionTracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: nil}

	pipelineCtx, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UnrackedOneParamsOther(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

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
	subscriptionTracker.SeenEndpointDescs[&ed] = pc

	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p1}

	pipelineCtx, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UntrackedBothParams(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

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
	subscriptionTracker.SeenEndpointDescs[&ed] = pc

	p2 := make(map[string]string)
	p2["param"] = "value2"
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p2}

	pipelineCtx, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}

func TestSubscriptionTracker_GetActivePipelineFromEndpointDesc_UntrackedBothParamsDiffKeys(t *testing.T) {
	subscriptionTracker := SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)}

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
	subscriptionTracker.SeenEndpointDescs[&ed] = pc

	p2 := make(map[string]string)
	p2["param2"] = "value"
	ed2 := EndpointDesc{Endpoint: "/firsttest", Params: &p2}

	pipelineCtx, err := subscriptionTracker.GetActivePipelineFromEndpointDesc(ed2)

	if err == nil {
		t.Errorf("should return an error")
	}

	if pipelineCtx == pc {
		t.Errorf("params in second desc don't match, so shouldn't return a pipeline")
	}

	cf()
}

//func TestSubscriptionTracker_TrackEndpointDesc_ExistNoParams(t *testing.T) {
//	subscriptionTracker := SubscriptionTracker{make(map[*PipeContext]*EndpointDesc)}
//
//	ctrl := gomock.NewController(t)
//
//	// Assert that Bar() is invoked.
//	defer ctrl.Finish()
//
//	m := pubsub_mock.NewMockPubSub(ctrl)
//	c, cf := context.WithCancel(context.Background())
//	pc := NewPipeContextFromContext(config.RouteConfig{}, make(gin.Params, 0), m, "/test", "testid", c)
//	ed := EndpointDesc{Endpoint: "/firsttest", Params: nil}
//	subscriptionTracker.TrackEndpointDesc(pc, ed)
//	assert.Equal(t, subscriptionTracker.SeenEndpointDescs[pc], &ed)
//	cf()
//}
