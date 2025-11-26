package httpclientv2

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	openapi_types "github.com/oapi-codegen/runtime/types"

	http_client "gitlab.noway/gen/http/profile_v2/client"
)

func (c *Client) Create(ctx context.Context, name string, age int, email, phone string) (uuid.UUID, error) {
	input := http_client.CreateProfileInput{
		Name:  name,
		Age:   age,
		Email: openapi_types.Email(email),
		Phone: phone,
	}

	output, err := c.client.CreateProfileWithResponse(ctx, input)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create profile: %w", err)
	}

	if output.StatusCode() != http.StatusOK {
		return uuid.Nil, fmt.Errorf("create profile: %w", errors.New(output.JSON400.Error))
	}

	return output.JSON200.ID, nil
}
