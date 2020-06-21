package server

import (
	"errors"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"

	"github.com/foxcool/homework/7-prometheus/app/api"
	"github.com/foxcool/homework/7-prometheus/app/logger"
	"github.com/foxcool/homework/7-prometheus/app/storage"
)

// POST /auth
func (s *Server) Auth(ctx echo.Context) error {
	var params api.AuthJSONRequestBody //nolint:govet
	if err := ctx.Bind(&params); err != nil {
		return s.error(ctx, http.StatusBadRequest, err)
	}

	responseInternalServerError := func(err error) error {
		s.authLog.Info(logger.AUTHUSERLOGIN, params.Login, logger.AUTHFAIL)
		panic(err)
	}

	responseUnauthorized := func(err error) error {
		s.authLog.Info(logger.AUTHUSERLOGIN, params.Login, logger.AUTHFAIL)
		return s.error(ctx, http.StatusUnauthorized, err)
	}

	responseSuccess := func(userID string) error {
		s.authLog.Info(logger.AUTHUSERLOGIN, params.Login, logger.AUTHSUCCESS)

		return ctx.JSON(http.StatusOK, api.Auth200{
			Success: api.Success{
				Base: api.Base{
					Version: s.Version,
				},
				Message: SuccessMessage,
			},
			Data: &[]struct {
				UserID string `json:"userID"`
			}{{
				UserID: userID,
			}},
		})
	}

	var input storage.User
	if strfmt.IsEmail(params.Login) {
		input.Email = &params.Login
	} else {
		input.Mobile = &params.Login
	}

	user, err := s.storage.GetUser(input)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return responseUnauthorized(err)
		}

		return responseInternalServerError(err)
	}

	newHash, err := VerifyPassword(params.Password, *user.PasswordHash)
	if err != nil {
		return responseUnauthorized(err)
	}

	if newHash != "" {
		user, err = s.storage.UpdateUser(*user.ID, storage.User{
			PasswordHash: &newHash,
		})
		if err != nil {
			return responseInternalServerError(err)
		}
	}



	return responseSuccess(*user.ID)
}
