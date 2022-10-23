package pipeline

import (
	"fmt"
)

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{
		tracker:             &SubscriptionTracker{make(map[*EndpointDesc]*PipeContext)},
		SubscribeRequests:   make(chan RouteConfigWithParams, 8),
		UnsubscribeRequests: make(chan RouteConfigWithParams, 8),
	}
}

func (sh *SubscriptionHandler) Track(ctx *PipeContext, desc *EndpointDesc) error {
	if _, err := sh.GetPipeCtx(*desc); err == nil {
		return fmt.Errorf("endpoint desc is already being handled by a pipeline")

	}

	// No endpoint names matched
	sh.tracker.SeenEndpointDescs[desc] = ctx
	return nil
}

func (sh *SubscriptionHandler) GetPipeCtx(desc EndpointDesc) (*PipeContext, error) {
	for endpointDesc, pipeCtx := range sh.tracker.SeenEndpointDescs {
		if endpointDescsMatch(*endpointDesc, desc) {
			return pipeCtx, nil
		}
	}
	return nil, fmt.Errorf("specified endpoint desc doesn't exist")
}

func (t *SubscriptionHandler) CancelAndDeleteEntry(p *PipeContext) error {
	for ed, pc := range t.tracker.SeenEndpointDescs {
		if pc == p {
			pc.CancelFunc()
			delete(t.tracker.SeenEndpointDescs, ed)
			return nil
		}
	}
	return fmt.Errorf("not currently tracking given pipeline context")
}

func endpointDescsMatch(e1 EndpointDesc, e2 EndpointDesc) bool {
	if e1.Endpoint == e2.Endpoint {
		// one specifies params and one doesn't
		if (e1.Params != nil && e2.Params == nil) || (e1.Params == nil && e2.Params != nil) {
			return false
		}

		// If they both don't, then they're the same
		if e1.Params == nil && e2.Params == nil {
			return true
		}

		// Both have params, see if they're the same
		for k, v := range *e1.Params {
			v2, ok := (*e2.Params)[k]
			// Different!
			if !ok {
				return false
			}

			if v2 != v {
				return false
			}
		}
		// They're the same
		return true
	}
	return false
}
