package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"password-validator/adapter/handler"
	"password-validator/adapter/response"
	"password-validator/core/usecase"
	"password-validator/core/usecase/input"

	"go.opentelemetry.io/otel/codes"

	"github.com/itau-corp/itau-jw1-dep-golibs-gotel/logger"
	oteltrace "github.com/itau-corp/itau-jw1-dep-golibs-gotel/otel/trace"
	"github.com/itau-corp/itau-jw1-dep-golibs-gotel/otel/utils"
)

type ValidatePasswordController struct {
	validatePasswordUseCase usecase.ValidatePasswordUseCase
}

func NewValidatePasswordController(
	validatePasswordUseCase usecase.ValidatePasswordUseCase,
) ValidatePasswordController {
	return ValidatePasswordController{
		validatePasswordUseCase: validatePasswordUseCase,
	}
}

func (c ValidatePasswordController) Execute(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	log.Info("ValidatePasswordController controller initialized")
	newCtx, span := oteltrace.NewSpan(r.Context(), "password-validator", "password-span")
	defer span.End()

	jsonBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error("Error reading request body", err)
		span.SetStatus(codes.Error, "ValidatePasswordController Error")
		span.RecordError(err)
		response.NewError(err, http.StatusBadRequest, nil).Send(w)
		return
	}

	var i input.PasswordInput
	if err := json.Unmarshal(jsonBody, &i); err != nil {
		log.Error("error unmarshal password input", err)
		handler.HandleErrors(w, err, nil)
		return
	}

	span.SetAttributes(utils.StringAttribute("password", i.Password))

	output, err := c.validatePasswordUseCase.Execute(newCtx, i)
	if err != nil {
		span.SetStatus(codes.Error, "ValidatePasswordController Error")
		span.RecordError(err)
		handler.HandleErrors(w, err, output)
		return
	}

	span.AddEvent("Finished ValidatePasswordController execution")
	span.SetStatus(codes.Ok, "ValidatePasswordController execution finished with success")
	response.NewSuccess(output, http.StatusOK).Send(w)
}
