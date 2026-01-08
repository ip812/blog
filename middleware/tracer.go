package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

func TraceIDHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		if span != nil {
			sc := span.SpanContext()
			if sc.IsValid() {
				w.Header().Set("X-Trace-Id", sc.TraceID().String())
			}
		}
		next.ServeHTTP(w, r)
	})
}

