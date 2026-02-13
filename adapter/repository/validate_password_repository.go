package repository

import (
	"context"
	"password-validator/core/domain/password"
	_errors "password-validator/core/errors"
)

type PasswordRepository struct {
	storage []password.Password
}

func NewPasswordRepository() *PasswordRepository {
	return &PasswordRepository{
		storage: []password.Password{},
	}
}

func (r *PasswordRepository) Save(ctx context.Context, p *password.Password) error {
	r.storage = append(r.storage, *p)
	return nil
}

func (r *PasswordRepository) FindById(ctx context.Context, pa string) (*password.Password, error) {
	for _, v := range r.storage {
		if v.Password() == pa {
			return &v, nil
		}
	}
	return &password.Password{}, _errors.NotFoundError{
		Entity: "Password",
		ID:     pa,
	}
}
