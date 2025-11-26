package usecase_test

import (
	"context"
	"testing"

	kafkalib "github.com/segmentio/kafka-go"

	"github.com/stretchr/testify/require"
	"gitlab.noway/internal/usecase"

	"gitlab.noway/internal/usecase/mocks"
	"gitlab.noway/pkg/otel"
	"gitlab.noway/pkg/transaction"
)

func Test_OutboxReadAndProduce_Success(t *testing.T) {
	otel.SilentModeInit()
	transaction.IsUnitTest = true

	// Данные для поведения
	msgs := []kafkalib.Message{{}, {}, {}}

	// Настраиваем поведение Postgres
	postgres := new(mocks.Postgres)
	postgres.On("ReadOutboxKafka", Any, Any).Return(msgs, nil)
	defer postgres.AssertCalled(t, "ReadOutboxKafka", Any, Any)

	//	Настройка поведения Kafka
	kafka := new(mocks.Kafka)
	kafka.On("Produce", Any, Any).Return(nil)
	defer kafka.AssertCalled(t, "Produce", Any, Any)

	// Собираем UseCase
	u := usecase.New(postgres, nil, kafka, nil)

	{ // Сам тест
		lenMessages, err := u.OutboxReadAndProduce(context.Background(), 10)
		require.NoError(t, err)
		require.Equal(t, len(msgs), lenMessages)
	}
}
