package dto

import (
	"gitlab.noway/internal/domain"
)

type UpdateProfileInput struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

func (i UpdateProfileInput) Validate() error {
	if i.Name == nil && i.Age == nil && i.Email == nil && i.Phone == nil {
		return domain.ErrAllFieldsForUpdateEmpty
	}

	return nil
}
