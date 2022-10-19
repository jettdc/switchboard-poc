package pipeline

import "net/http"

type MiddlewarePlugin interface {
	Process(r *http.Request) error
}

type EnrichmentPlugin interface {
	Process(payload string) (string, error)
}
