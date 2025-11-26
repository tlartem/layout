//go:build http_v1

package test

import (
	"github.com/google/uuid"
)

func (s *Suite) Test_Kafka() {
	msg, err := s.kafkaReader.ReadMessage(ctx)
	s.NoError(err)

	s.NoError(uuid.Validate(string(msg.Key)))
	s.NoError(uuid.Validate(string(msg.Value)))
}
