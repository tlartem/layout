package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "gitlab.noway/gen/grpc/profile_v1"
	"gitlab.noway/internal/dto"
	"gitlab.noway/internal/dto/baggage"
)

func (h Handlers) DeleteProfile(ctx context.Context, i *pb.DeleteProfileInput) (*emptypb.Empty, error) {
	input := dto.DeleteProfileInput{
		ID: i.GetId(),
	}

	baggage.PutProfileID(ctx, input.ID)

	err := h.usecase.DeleteProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
