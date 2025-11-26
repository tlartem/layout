package logger

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"gitlab.noway/internal/dto/baggage"
)

func Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	bag := &baggage.Baggage{}
	ctx = baggage.WithContext(ctx, bag)

	event := log.Info()

	resp, err := handler(ctx, req)
	if err != nil {
		event = log.Error().Err(bag.Err)
	}

	event.
		Str("profile_id", bag.ProfileID).
		Str("proto", "grpc").
		Str("code", status.Code(err).String()).
		Str("method", info.FullMethod).
		Send()

	return resp, err
}
