package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int32     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Age       int32     `json:"age" db:"age" validate:"required,min=1,max=150"`
	Role      string    `json:"role" db:"role" validate:"required,oneof=admin user moderator"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserParams struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Age   int32  `json:"age" validate:"required,min=1,max=150"`
	Role  string `json:"role" validate:"required,oneof=admin user moderator"`
}

type UpdateUserParams struct {
	ID    int32  `json:"id" validate:"required,min=1"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Age   int32  `json:"age" validate:"required,min=1,max=150"`
	Role  string `json:"role" validate:"required,oneof=admin user moderator"`
}

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// Validation functions
func (u *CreateUserParams) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UpdateUserParams) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
