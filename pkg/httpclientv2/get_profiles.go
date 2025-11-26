package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	http_client "gitlab.noway/gen/http/profile_v2/client"
)

func (c *Client) GetProfiles(ctx context.Context, sort, order string, offset, limit int,
) ([]http_client.GetProfileOutput, error) {
	params := http_client.GetProfilesParams{
		Sort:   sort,
		Order:  &order,
		Offset: &offset,
		Limit:  &limit,
	}

	output, err := c.client.GetProfilesWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("GetProfilesWithResponse: %w", err)
	}

	if output.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if output.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request failed: status: %s, body:%s", output.Status(), output.Body)
	}

	return *output.JSON200, nil
}
