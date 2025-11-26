package domain

import "github.com/google/uuid"

type Property struct {
	ProfileID uuid.UUID
	Tags      []string
}

func NewProperty(profileID uuid.UUID, tags []string) Property {
	return Property{
		ProfileID: profileID,
		Tags:      tags,
	}
}
