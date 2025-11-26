package v2

import (
	"context"

	http_server "gitlab.noway/gen/http/profile_v2/server"
	"gitlab.noway/internal/dto"
)

func (h *Handlers) UpdateProfile(ctx context.Context, request http_server.UpdateProfileRequestObject,
) (http_server.UpdateProfileResponseObject, error) {
	input := dto.UpdateProfileInput{
		ID:    request.Body.ID.String(),
		Name:  request.Body.Name,
		Age:   request.Body.Age,
		Email: request.Body.Email,
		Phone: request.Body.Phone,
	}

	err := h.usecase.UpdateProfile(ctx, input)
	if err != nil {
		return http_server.UpdateProfile400JSONResponse{Error: err.Error()}, err
	}

	return http_server.UpdateProfile204Response{}, err
}
