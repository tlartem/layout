package router

import (
	"context"
	"strings"

	"github.com/go-chi/chi/v5"
)

func ExtractPath(ctx context.Context) string {
	chiCtx := chi.RouteContext(ctx)
	path := strings.Join(chiCtx.RoutePatterns, "")
	path = strings.ReplaceAll(path, "/*/", "/")

	return path
}
