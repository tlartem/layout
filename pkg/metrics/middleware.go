package metrics

import (
	"net/http"
	"time"

	"gitlab.noway/pkg/router"
)

func NewMiddleware(metrics *HTTPServer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			ww := router.WriterWrapper(w)
			next.ServeHTTP(ww, r)

			method := r.Method + " " + router.ExtractPath(r.Context())

			// Metrics
			metrics.Duration(method, now)
			metrics.TotalInc(method, ww.Code())
		}

		return http.HandlerFunc(fn)
	}
}
