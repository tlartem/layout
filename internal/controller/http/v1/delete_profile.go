package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.noway/internal/dto"
	"gitlab.noway/internal/dto/baggage"
	"gitlab.noway/pkg/render"
)

func (h *Handlers) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.DeleteProfileInput{
		ID: chi.URLParam(r, "id"),
	}

	baggage.PutProfileID(ctx, input.ID)

	err := h.usecase.DeleteProfile(ctx, input)
	if err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "request failed")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
