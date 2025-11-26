package dto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.noway/internal/dto"
)

func Test_GetProfilesInput_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   dto.GetProfilesInput
		wantErr bool
	}{
		{
			name: "valid input id asc",
			input: dto.GetProfilesInput{
				Sort:   "id",
				Order:  "asc",
				Offset: 0,
				Limit:  10,
			},
			wantErr: false,
		},
		{
			name: "valid input id asc case insensitive",
			input: dto.GetProfilesInput{
				Sort:   "ID",
				Order:  "ASC",
				Offset: 0,
				Limit:  10,
			},
			wantErr: false,
		},
		{
			name: "valid input name desc",
			input: dto.GetProfilesInput{
				Sort:   "name",
				Order:  "desc",
				Offset: 0,
				Limit:  10,
			},
			wantErr: false,
		},
		{
			name: "only required field",
			input: dto.GetProfilesInput{
				Sort:   "name",
				Order:  "",
				Offset: 0,
				Limit:  0,
			},
			wantErr: false,
		},
		{
			name: "invalid sort field",
			input: dto.GetProfilesInput{
				Sort:   "invalid",
				Order:  "asc",
				Offset: 0,
				Limit:  10,
			},
			wantErr: true,
		},
		{
			name: "required sort field",
			input: dto.GetProfilesInput{
				Sort:   "",
				Order:  "",
				Offset: 0,
				Limit:  10,
			},
			wantErr: true,
		},
		{
			name: "invalid order field",
			input: dto.GetProfilesInput{
				Sort:   "id",
				Order:  "invalid",
				Offset: 0,
				Limit:  10,
			},
			wantErr: true,
		},
		{
			name: "negative offset",
			input: dto.GetProfilesInput{
				Sort:   "id",
				Order:  "asc",
				Offset: -1,
				Limit:  10,
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			input: dto.GetProfilesInput{
				Sort:   "id",
				Order:  "asc",
				Offset: 0,
				Limit:  -1,
			},
			wantErr: true,
		},
		{
			name: "limit out of range",
			input: dto.GetProfilesInput{
				Sort:   "id",
				Order:  "asc",
				Offset: 0,
				Limit:  101,
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
