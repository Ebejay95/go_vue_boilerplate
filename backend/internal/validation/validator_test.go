package validation

import (
	"testing"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required,min=2,max=50"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,min=1,max=150"`
	Role  string `validate:"required,oneof=admin user moderator"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name      string
		input     TestStruct
		wantError bool
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
				Role:  "user",
			},
			wantError: false,
		},
		{
			name: "invalid email",
			input: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   30,
				Role:  "user",
			},
			wantError: true,
		},
		{
			name: "missing required field",
			input: TestStruct{
				Email: "john@example.com",
				Age:   30,
				Role:  "user",
			},
			wantError: true,
		},
		{
			name: "invalid role",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
				Role:  "invalid_role",
			},
			wantError: true,
		},
		{
			name: "age too young",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   0,
				Role:  "user",
			},
			wantError: true,
		},
		{
			name: "age too old",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   200,
				Role:  "user",
			},
			wantError: true,
		},
		{
			name: "name too short",
			input: TestStruct{
				Name:  "J",
				Email: "john@example.com",
				Age:   30,
				Role:  "user",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		valid    bool
	}{
		{"valid username", "john_doe123", true},
		{"valid short username", "joe", true},
		{"valid long username", "very_long_username_that_is_ok", true},
		{"invalid too short", "jo", false},
		{"invalid too long", "this_username_is_way_too_long_to_be_valid", false},
		{"invalid special chars", "john-doe", false},
		{"invalid spaces", "john doe", false},
		{"invalid symbols", "john@doe", false},
		{"empty username", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to create a validator instance to test the custom validation
			validate := validator.New()
			validate.RegisterValidation("username", validateUsername)

			type TestUser struct {
				Username string `validate:"username"`
			}

			testUser := TestUser{Username: tt.username}
			err := validate.Struct(testUser)

			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
