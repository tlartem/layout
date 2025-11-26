package usecase_test

import (
	"context"
	"testing"
	"time"

	"gitlab.noway/pkg/otel"

	"gitlab.noway/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.noway/internal/domain"
	"gitlab.noway/internal/dto"
	"gitlab.noway/internal/usecase/mocks"
)

// Покрытие кода тестами
// go test -coverprofile=cover.out ./internal/usecase && go tool cover -html=cover.out

const Any = mock.Anything

func Test_GetProfile_Success(t *testing.T) {
	otel.SilentModeInit()

	// Данные для поведения
	id := uuid.New()
	profile := domain.Profile{ID: id}

	// Настраиваем поведение Postgres
	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", Any, id).Return(profile, nil)
	defer postgres.AssertCalled(t, "GetProfile", Any, id)

	// Собираем UseCase
	u := usecase.New(postgres, nil, nil, nil)

	{ // Сам тест
		input := dto.GetProfileInput{ID: id.String()}
		output := dto.GetProfileOutput{Profile: profile}

		actual, err := u.GetProfile(context.Background(), input)
		require.NoError(t, err)
		require.Equal(t, output, actual)
	}
}

func Test_GetProfile_InvalidUUID(t *testing.T) {
	otel.SilentModeInit()

	// Собираем UseCase
	u := usecase.New(nil, nil, nil, nil)

	{ // Сам тест
		input := dto.GetProfileInput{ID: "invalid-uuid"}
		output := dto.GetProfileOutput{}

		actual, err := u.GetProfile(context.Background(), input)
		require.Error(t, err)
		require.Equal(t, output, actual)
		require.Equal(t, domain.ErrUUIDInvalid, err)
	}
}

func Test_GetProfile_NotFound(t *testing.T) {
	otel.SilentModeInit()

	// Данные для поведения
	id := uuid.New()

	// Настраиваем поведение Postgres
	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", Any, id).Return(domain.Profile{}, domain.ErrNotFound)
	defer postgres.AssertCalled(t, "GetProfile", Any, id)

	// Собираем UseCase
	u := usecase.New(postgres, nil, nil, nil)

	{ // Сам тест
		input := dto.GetProfileInput{ID: id.String()}
		output := dto.GetProfileOutput{}

		actual, err := u.GetProfile(context.Background(), input)
		require.Error(t, err)
		require.Equal(t, output, actual)
		require.ErrorIs(t, err, domain.ErrNotFound)
	}
}

func Test_GetProfile_DeletedProfile(t *testing.T) {
	otel.SilentModeInit()

	// Данные для поведения
	id := uuid.New()
	profile := domain.Profile{ID: id, DeletedAt: time.Now()}

	// Настраиваем поведение Postgres
	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", Any, id).Return(profile, nil)
	defer postgres.AssertCalled(t, "GetProfile", Any, id)

	// Собираем UseCase
	u := usecase.New(postgres, nil, nil, nil)

	{ // Сам тест
		input := dto.GetProfileInput{ID: id.String()}
		output := dto.GetProfileOutput{}

		actual, err := u.GetProfile(context.Background(), input)
		require.Error(t, err)
		require.Equal(t, output, actual)
		require.ErrorIs(t, err, domain.ErrNotFound)
	}
}
