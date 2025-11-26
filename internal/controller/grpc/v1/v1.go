package v1

import (
	pb "gitlab.noway/gen/grpc/profile_v1"
	"gitlab.noway/internal/usecase"
)

type Handlers struct {
	pb.UnimplementedProfileV1Server
	usecase *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handlers {
	return &Handlers{
		usecase: uc,
	}
}
