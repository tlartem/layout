package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "gitlab.noway/gen/grpc/profile_v1"
	"gitlab.noway/internal/dto"
	"gitlab.noway/internal/dto/baggage"
)

func (h Handlers) CreateProfile(ctx context.Context, i *pb.CreateProfileInput) (*pb.CreateProfileOutput, error) {
	input := dto.CreateProfileInput{
		Name:  i.GetName(),
		Age:   int(i.GetAge()),
		Email: i.GetEmail(),
		Phone: i.GetPhone(),
	}

	output, err := h.usecase.CreateProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	baggage.PutProfileID(ctx, output.ID.String())

	return &pb.CreateProfileOutput{
		Id: output.ID.String(),
	}, nil
}
