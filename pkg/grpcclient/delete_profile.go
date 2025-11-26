package grpcclient

import (
	"context"
	"fmt"

	pb "gitlab.noway/gen/grpc/profile_v1"
)

func (c *Client) Delete(ctx context.Context, id string) error {
	input := &pb.DeleteProfileInput{
		Id: id,
	}

	_, err := c.client.DeleteProfile(ctx, input)
	if err != nil {
		return fmt.Errorf("client.DeleteProfile: %w", err)
	}

	return nil
}
