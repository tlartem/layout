package v2

import (
	"context"

	http_server "gitlab.noway/gen/http/profile_v2/server"
	"gitlab.noway/internal/dto"
)

func (h *Handlers) DeleteProfileByID(ctx context.Context, request http_server.DeleteProfileByIDRequestObject,
) (http_server.DeleteProfileByIDResponseObject, error) {
	input := dto.DeleteProfileInput{
		ID: request.ID.String(),
	}

	err := h.usecase.DeleteProfile(ctx, input)
	if err != nil {
		return http_server.DeleteProfileByID400JSONResponse{Error: err.Error()}, nil //nolint:nilerr
	}

	return http_server.DeleteProfileByID204Response{}, nil
}
