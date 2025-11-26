package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	http_server "gitlab.noway/gen/http/profile_v2/server"
	ver1 "gitlab.noway/internal/controller/http/v1"
	ver2 "gitlab.noway/internal/controller/http/v2"
	"gitlab.noway/internal/usecase"
	"gitlab.noway/pkg/logger"
	"gitlab.noway/pkg/metrics"
	"gitlab.noway/pkg/otel"
)

func ProfileRouter(r *chi.Mux, uc *usecase.UseCase, m *metrics.HTTPServer) {
	v1 := ver1.New(uc)
	v2 := ver2.New(uc)

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/noway/layout/api", func(r chi.Router) {
		r.Use(logger.Middleware)
		r.Use(metrics.NewMiddleware(m))
		r.Use(otel.Middleware)

		r.Route("/v1", func(r chi.Router) {
			r.Post("/profile", v1.CreateProfile)
			r.Put("/profile", v1.UpdateProfile)
			r.Get("/profile/{id}", v1.GetProfile)
			r.Get("/profiles", v1.GetProfiles)
			r.Delete("/profile/{id}", v1.DeleteProfile)
		})

		r.Route("/v2", func(r chi.Router) {
			mux := http_server.NewStrictHandler(v2, []http_server.StrictMiddlewareFunc{})
			http_server.HandlerFromMux(mux, r)
		})
	})
}
