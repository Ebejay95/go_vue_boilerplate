package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUserParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  CreateUserParams
		wantErr bool
	}{
		{
			name: "valid user params",
			params: CreateUserParams{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
				Role:  "user",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			params: CreateUserParams{
				Email: "john@example.com",
				Age:   25,
				Role:  "user",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			params: CreateUserParams{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   25,
				Role:  "user",
			},
			wantErr: true,
		},
		{
			name: "invalid age",
			params: CreateUserParams{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   0,
				Role:  "user",
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			params: CreateUserParams{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
				Role:  "invalid_role",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUserParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  UpdateUserParams
		wantErr bool
	}{
		{
			name: "valid update params",
			params: UpdateUserParams{
				ID:    1,
				Name:  "Jane Doe",
				Email: "jane@example.com",
				Age:   30,
				Role:  "admin",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			params: UpdateUserParams{
				Name:  "Jane Doe",
				Email: "jane@example.com",
				Age:   30,
				Role:  "admin",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
