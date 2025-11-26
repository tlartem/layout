package usecase

import (
	"context"

	"gitlab.noway/pkg/otel/tracer"
)

//nolint:gocritic
func (u *UseCase) SomeWork(ctx context.Context) error {
	_, span := tracer.Start(ctx, "usecase SomeWork")
	defer span.End()

	// log.Info().Msg("SomeWork called")

	// Пример вызова клиента
	// p, err := u.profile.GetProfile(ctx, "8638341a-b68a-4291-84ee-94b147afeff9")
	// if err != nil {
	//	return fmt.Errorf("SomeWork: %w", err)
	//}
	//
	// fmt.Println(p)

	return nil
}
