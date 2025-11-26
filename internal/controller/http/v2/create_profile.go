package v2

import (
	"context"

	http_server "gitlab.noway/gen/http/profile_v2/server"
	"gitlab.noway/internal/dto"
)

func (h *Handlers) CreateProfile(ctx context.Context, request http_server.CreateProfileRequestObject,
) (http_server.CreateProfileResponseObject, error) {
	input := dto.CreateProfileInput{
		Name:  request.Body.Name,
		Age:   request.Body.Age,
		Email: string(request.Body.Email),
		Phone: request.Body.Phone,
	}

	output, err := h.usecase.CreateProfile(ctx, input)
	if err != nil {
		return http_server.CreateProfile400JSONResponse{Error: err.Error()}, nil //nolint:nilerr
	}

	return http_server.CreateProfile200JSONResponse{
		ID: output.ID,
	}, nil
}
