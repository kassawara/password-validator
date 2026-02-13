package repository

import (
	"context"
	"password-validator/core/domain/password"
)

type PasswordRepository interface {
	Save(context.Context, *password.Password) error
	FindById(context.Context, string) (*password.Password, error)
}
