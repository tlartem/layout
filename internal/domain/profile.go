package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Name string

type Age int

type Profile struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	Name      Name      `json:"name"       validate:"required,min=3,max=64"`
	Age       Age       `json:"age"        validate:"required,min=18,max=120"`
	Status    Status    `json:"status"`
	Verified  bool      `json:"verified"`
	Contacts  Contacts  `json:"contacts"`
}

type Contacts struct {
	Email string `json:"email" validate:"email"`
	Phone string `json:"phone" validate:"e164"`
}

var validate = validator.New(validator.WithRequiredStructEnabled()) //nolint:gochecknoglobals

func NewProfile(name string, age int, email, phone string) (Profile, error) {
	p := Profile{
		ID:       uuid.New(),
		Name:     Name(name),
		Age:      Age(age),
		Status:   Pending,
		Verified: false,
		Contacts: Contacts{
			Email: email,
			Phone: phone,
		},
	}

	if err := p.Validate(); err != nil {
		return Profile{}, fmt.Errorf("p.Validate: %w", err)
	}

	return p, nil
}

func (p Profile) Validate() error {
	err := validate.Struct(p)
	if err != nil {
		return fmt.Errorf("validate.Struct: %w", err)
	}

	return nil
}

func (p Profile) IsDeleted() bool {
	return !p.DeletedAt.IsZero()
}
