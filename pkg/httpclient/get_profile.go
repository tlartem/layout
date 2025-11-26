package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Profile struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Status    int    `json:"status"`
	Verified  bool   `json:"verified"`
	Contacts  struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contacts"`
}

func (c *Client) GetProfile(ctx context.Context, id string) (Profile, error) {
	const getProfile = "noway/layout/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, getProfile, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return Profile{}, fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Profile{}, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Profile{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return Profile{}, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return Profile{}, fmt.Errorf("request failed: status: %s, body:%s", resp.Status, body)
	}

	var profile Profile

	if err = json.Unmarshal(body, &profile); err != nil {
		return Profile{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return profile, nil
}
