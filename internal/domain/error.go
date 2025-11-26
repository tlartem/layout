package domain

import (
	"errors"
)

var (
	ErrNotFound                = errors.New("not found")
	ErrAllFieldsForUpdateEmpty = errors.New("all fields for update are empty")
	ErrUUIDInvalid             = errors.New("uuid is invalid")
	ErrEmptyTopic              = errors.New("topic is empty")
)
