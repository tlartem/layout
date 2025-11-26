package v2

import (
	"context"
	"errors"

	http_server "gitlab.noway/gen/http/profile_v2/server"
	"gitlab.noway/internal/domain"
	"gitlab.noway/internal/dto"
)

func (h *Handlers) GetProfiles(ctx context.Context, request http_server.GetProfilesRequestObject,
) (http_server.GetProfilesResponseObject, error) {
	input := dto.GetProfilesInput{
		Sort: request.Params.Sort,
	}

	if request.Params.Order != nil {
		input.Order = *request.Params.Order
	}

	if request.Params.Offset != nil {
		input.Offset = *request.Params.Offset
	}

	if request.Params.Limit != nil {
		input.Limit = *request.Params.Limit
	}

	output, err := h.usecase.GetProfiles(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return http_server.GetProfiles404JSONResponse{Error: err.Error()}, nil

		default:
			return http_server.GetProfiles400JSONResponse{Error: err.Error()}, nil
		}
	}

	profiles := make(http_server.GetProfiles200JSONResponse, 0, len(output.Profiles))

	for _, profile := range output.Profiles {
		var p http_server.GetProfileOutput

		p.ID = profile.ID
		p.CreatedAt = profile.CreatedAt
		p.UpdatedAt = profile.UpdatedAt
		p.Name = string(profile.Name)
		p.Age = int(profile.Age)
		p.Status = int(profile.Status)
		p.Verified = profile.Verified
		p.Contacts.Email = profile.Contacts.Email
		p.Contacts.Phone = profile.Contacts.Phone

		profiles = append(profiles, p)
	}

	return profiles, nil
}
