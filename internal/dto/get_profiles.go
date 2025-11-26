package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"gitlab.noway/internal/domain"
)

type GetProfilesOutput struct {
	Profiles []domain.Profile `json:"profiles"`
}

type GetProfilesInput struct {
	Sort   string `validate:"oneofci=id name"`
	Order  string `validate:"oneofci='' asc desc"`
	Offset int    `validate:"min=0"`
	Limit  int    `validate:"min=0,max=100"`
}

var validate = validator.New(validator.WithRequiredStructEnabled()) //nolint:gochecknoglobals

func (i GetProfilesInput) Validate() error {
	err := validate.Struct(i)
	if err != nil {
		return fmt.Errorf("validate.Struct: %w", err)
	}

	return nil
}
