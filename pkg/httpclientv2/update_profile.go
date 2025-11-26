package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	http_client "gitlab.noway/gen/http/profile_v2/client"
)

func (c *Client) Update(ctx context.Context, id string, name *string, age *int, email, phone *string) error {
	input := http_client.UpdateProfileInput{
		ID:    uuid.MustParse(id),
		Name:  name,
		Age:   age,
		Email: email,
		Phone: phone,
	}

	output, err := c.client.UpdateProfileWithResponse(ctx, input)
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	if output.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("request failed: status: %s, body:%s", output.Status(), output.Body)
	}

	return nil
}
