package middleware

import (
	"net/http"

	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
)

// NetHTTP is middleware for HTTP handler
func NetHTTP(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			span := &trace.Span{}
			defer span.End()

			ctx := r.Context()
			hf := propagation.HTTPFormat{}
			if sc, ok := hf.SpanContextFromRequest(r); ok {
				ctx, span = trace.StartSpanWithRemoteParent(ctx, label, sc)
			} else {
				ctx, span = trace.StartSpan(ctx, label)
			}

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
