package types

import (
	"context"
	"time"

	"github.com/jayden1905/abundance/cmd/pkg/database"
)

type User struct {
	ID           int32     `json:"id"`
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	Subscription string    `json:"subscription"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	IsVerified   bool      `json:"is_verify"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserStore interface {
	GetUsersPaginated(page int32, pageSize int32) ([]*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int32) (*User, error)
	GetUserRoleByID(id int32) (string, error)
	CreateUser(ctx context.Context, user *database.User) error
	CreateSuperUser(ctx context.Context, user *User) error
	UpdateUserToSuperUser(ctx context.Context, id int32) error
	UpdateUserToNormalUser(ctx context.Context, id int32) error
	UpdateUserVerification(ctx context.Context, id int32) error
	DeleteUserByID(ctx context.Context, id int32) error
}

type RegisterUserPayload struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=3,max=20"`
	Role         string `json:"role" validate:"required"`
	Subscription string `json:"subscription" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserInformationPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type ResendVerificationEmailPayload struct {
	Email string `json:"email" validate:"required,email"`
}
