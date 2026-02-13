package usecase

import (
	"context"
	"password-validator/core/domain/password"
	"password-validator/core/repository"
	"password-validator/core/usecase/input"
	"password-validator/core/usecase/output"
	"time"

	"github.com/itau-corp/itau-jw1-dep-golibs-gotel/logger"
)

type (
	ValidatePasswordUseCase interface {
		Execute(context.Context, input.PasswordInput) (output.PasswordOutput, error)
	}

	ValidatePasswordPresenter interface {
		Output(context.Context, *password.Password) output.PasswordOutput
	}

	validatePasswordUseCase struct {
		ctxTimeout time.Duration
		repository repository.PasswordRepository
		presenter  ValidatePasswordPresenter
	}
)

func NewValidatePasswordUseCase(
	ctxTimeout time.Duration,
	repository repository.PasswordRepository,
	presenter ValidatePasswordPresenter,
) ValidatePasswordUseCase {
	return &validatePasswordUseCase{
		ctxTimeout: ctxTimeout,
		repository: repository,
		presenter:  presenter,
	}
}

func (u validatePasswordUseCase) Execute(ctx context.Context, i input.PasswordInput) (output.PasswordOutput, error) {
	log := logger.FromContext(ctx).WithFields(logger.Field{"password": i.Password})
	log.Info("Validate password usecase initialized")

	p, err := password.New(
		password.WithPassword(i.Password),
	)
	if err != nil {
		return output.PasswordOutput{}, err
	}

	_ = u.repository.Save(ctx, p)

	log.Info("Validate password usecase finished")
	return u.presenter.Output(ctx, p), nil
}
