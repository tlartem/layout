package usecase_test

import (
	"context"
	"testing"

	"gitlab.noway/pkg/otel"

	"gitlab.noway/pkg/transaction"

	"github.com/stretchr/testify/require"
	"gitlab.noway/internal/dto"
	"gitlab.noway/internal/usecase"
	"gitlab.noway/internal/usecase/mocks"
)

func Test_CreateProfileV2_Success(t *testing.T) {
	otel.SilentModeInit()
	transaction.IsUnitTest = true

	// Настраиваем поведение Postgres
	postgres := new(mocks.Postgres)
	postgres.On("CreateProfile", Any, Any).Return(nil)
	postgres.On("CreateProperty", Any, Any).Return(nil)
	defer postgres.AssertCalled(t, "CreateProfile", Any, Any)
	defer postgres.AssertCalled(t, "CreateProperty", Any, Any)

	// Собираем UseCase
	u := usecase.New(postgres, nil, nil, nil)

	{ // Сам тест
		input := dto.CreateProfileInput{
			Name:  "John Doe",
			Age:   30,
			Email: "john.doe@example.com",
			Phone: "+1234567890",
		}

		actual, err := u.CreateProfileV2(context.Background(), input)
		require.NoError(t, err)
		require.NotEmpty(t, actual.ID)
	}
}
