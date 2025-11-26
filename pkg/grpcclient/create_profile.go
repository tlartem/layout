package grpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	pb "gitlab.noway/gen/grpc/profile_v1"
)

func (c *Client) Create(ctx context.Context, name string, age int, email, phone string) (uuid.UUID, error) {
	input := &pb.CreateProfileInput{
		Name:  name,
		Age:   int32(age), //nolint:gosec
		Email: email,
		Phone: phone,
	}

	resp, err := c.client.CreateProfile(ctx, input)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("c.client.CreateProfile: %w", err)
	}

	id, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("uuid.Parse: %w", err)
	}

	return id, nil
}
