package action

import (
	"errors"
	"github.com/gbooksoft/backend/internal/api/action/response"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gbooksoft/backend/internal/domain/password"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Authorize struct {
	service *authorization.Service
}

func NewAuthorize(service *authorization.Service) *Authorize {
	return &Authorize{
		service: service,
	}
}

type resetPasswordRequest struct {
	Hash     string `json:"token"`
	Password string `json:"password"`
}
type resetPasswordErrorResponse struct {
	PasswordErrors []string `json:"password_errors"`
}

func (a *Authorize) ResetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	err := a.service.ResetPassword(req.Hash, req.Password)

	if err != nil {
		var passErr *password.PasswordValidationResult
		if errors.As(err, &passErr) {
			response.Json(ctx.Writer).StatusCode(http.StatusBadRequest).Struct(resetPasswordErrorResponse{
				PasswordErrors: passErr.Errors,
			}).Write()
			return
		}

		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}
