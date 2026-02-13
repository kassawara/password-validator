package repository

import (
	"context"
	"password-validator/core/domain/password"

	"github.com/stretchr/testify/mock"
)

type PasswordRepositoryMock struct {
	mock.Mock
}

func (m *PasswordRepositoryMock) Save(ctx context.Context, p *password.Password) error {
	ret := m.Called(ctx, p)
	return ret.Error(0)
}

func (m *PasswordRepositoryMock) FindById(ctx context.Context, pa string) (*password.Password, error) {
	ret := m.Called(ctx, pa)
	return ret.Get(0).(*password.Password), ret.Error(1)
}
