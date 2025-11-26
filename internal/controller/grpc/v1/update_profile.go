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

func (h Handlers) UpdateProfile(ctx context.Context, i *pb.UpdateProfileInput) (*emptypb.Empty, error) {
	input := dto.UpdateProfileInput{
		ID:    i.GetId(),
		Name:  i.Name,
		Age:   parseAge(i.Age),
		Email: i.Email,
		Phone: i.Phone,
	}

	baggage.PutProfileID(ctx, input.ID)

	err := h.usecase.UpdateProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func parseAge(age *int32) *int {
	if age == nil {
		return nil
	}

	a := int(*age)

	return &a
}
