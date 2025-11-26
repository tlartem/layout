package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	http_client "gitlab.noway/gen/http/profile_v2/client"
)

func (c *Client) GetProfile(ctx context.Context, id string) (*http_client.GetProfileOutput, error) {
	output, err := c.client.GetProfileByIDWithResponse(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, fmt.Errorf("GetProfileByIdWithResponse: %w", err)
	}

	if output.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if output.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request failed: status: %s, body:%s", output.Status(), output.Body)
	}

	return output.JSON200, nil
}
