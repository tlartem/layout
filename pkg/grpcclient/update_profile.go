package grpcclient

import (
	"context"
	"fmt"

	pb "gitlab.noway/gen/grpc/profile_v1"
)

func (c *Client) Update(ctx context.Context, id string, name *string, age *int, email, phone *string) error {
	input := &pb.UpdateProfileInput{
		Id:    id,
		Name:  name,
		Age:   parseAge(age),
		Email: email,
		Phone: phone,
	}

	_, err := c.client.UpdateProfile(ctx, input)
	if err != nil {
		return fmt.Errorf("c.client.UpdateProfile: %w", err)
	}

	return nil
}

func parseAge(age *int) *int32 {
	if age == nil {
		return nil
	}

	a := int32(*age) //nolint:gosec

	return &a
}
