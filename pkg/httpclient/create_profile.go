package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) Create(ctx context.Context, name string, age int, email, phone string) (uuid.UUID, error) {
	const createProfile = "noway/layout/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s", c.host, createProfile)

	request := struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}{
		Name:  name,
		Age:   age,
		Email: email,
		Phone: phone,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return uuid.Nil, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return uuid.Nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return uuid.Nil, fmt.Errorf("request failed: status: %s, body:%s", resp.Status, body)
	}

	response := struct {
		ID uuid.UUID `json:"id"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return uuid.Nil, fmt.Errorf("json.Decode: %w", err)
	}

	return response.ID, nil
}
