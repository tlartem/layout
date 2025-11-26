package domain_test

import (
	"github.com/stretchr/testify/require"
	"gitlab.noway/internal/domain"
	"testing"
)

func TestNewProfile(t *testing.T) {
	cases := []struct {
		name    string
		age     int
		email   string
		phone   string
		wantErr bool
	}{
		{"Valid Profile", 25, "test@example.com", "+7123456789", false},
		{"Invalid Age Min", 17, "test@example.com", "+7123456789", true},
		{"Invalid Age Max", 121, "test@example.com", "+7123456789", true},
		{"Invalid Email", 25, "invalid-email", "+7123456789", true},
		{"Invalid Phone", 25, "test@example.com", "invalid-phone", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			profile, err := domain.NewProfile(c.name, c.age, c.email, c.phone)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.name, string(profile.Name))
				require.Equal(t, c.age, int(profile.Age))
				require.Equal(t, c.email, profile.Contacts.Email)
				require.Equal(t, c.phone, profile.Contacts.Phone)
				require.Equal(t, domain.Pending, profile.Status)
				require.False(t, profile.Verified)
			}
		})
	}
}
