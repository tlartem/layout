package logger

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func First(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	fmt.Println("First before", info.FullMethod)

	resp, err := handler(ctx, req)

	fmt.Println("First after", info.FullMethod)

	return resp, err
}

func Second(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	fmt.Println("Second before", info.FullMethod)

	resp, err := handler(ctx, req)

	fmt.Println("Second after", info.FullMethod)

	return resp, err
}
