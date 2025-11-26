package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.noway/internal/domain"
	"gitlab.noway/internal/dto"
	"gitlab.noway/internal/dto/baggage"
	"gitlab.noway/pkg/render"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.GetProfileInput{
		ID: chi.URLParam(r, "id"),
	}

	baggage.PutProfileID(ctx, input.ID)

	output, err := h.usecase.GetProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(ctx, w, err, http.StatusNotFound, "not found")

		default:
			render.Error(ctx, w, err, http.StatusBadRequest, "request failed")
		}

		return
	}

	render.JSON(w, output, http.StatusOK)
}
