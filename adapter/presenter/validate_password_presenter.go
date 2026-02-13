package presenter

import (
	"context"
	"password-validator/core/domain/password"
	"password-validator/core/usecase"
	"password-validator/core/usecase/output"
)

type validatePasswordPresenter struct{}

var _ usecase.ValidatePasswordPresenter = (*validatePasswordPresenter)(nil)

func NewValidatePasswordPresenter() usecase.ValidatePasswordPresenter {
	return &validatePasswordPresenter{}
}

func (p *validatePasswordPresenter) Output(ctx context.Context, password *password.Password) output.PasswordOutput {
	return output.PasswordOutput{
		IsValid: password.IsValid(),
	}
}
