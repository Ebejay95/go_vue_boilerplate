package models

import "time"

type User struct {
	ID        int32     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Age       int32     `json:"age" db:"age"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int32  `json:"age"`
	Role  string `json:"role"`
}

type UpdateUserParams struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int32  `json:"age"`
	Role  string `json:"role"`
}

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
