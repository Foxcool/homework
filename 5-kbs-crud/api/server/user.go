package server

import (
	"errors"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	apiErrors "github.com/foxcool/homework/5-k8s-crud/oapi-errors"

	"github.com/foxcool/homework/5-k8s-crud/api"
	"github.com/foxcool/homework/5-k8s-crud/logger"
	"github.com/foxcool/homework/5-k8s-crud/storage"
)

// POST /users
func (s *Server) PostUsers(ctx echo.Context) error {
	var params api.PostUsersJSONRequestBody
	if err := ctx.Bind(&params); err != nil {
		return s.error(ctx, http.StatusBadRequest, err)
	}

	responseInternalServerError := func(err error) error {
		s.authLog.Info(logger.USERCREATE, *params.Email, logger.USERFAIL)

		panic(err)
	}

	responseBadRequest := func(err error, validation map[string]interface{}) error {
		s.authLog.Info(logger.USERCREATE, *params.Email, logger.USERFAIL)

		response := apiErrors.ValidationErrorResponse{
			Response: apiErrors.Response{
				Version: s.Version,
				Message: apiErrors.ValidationErrorMessage,
			},
			Errors: &apiErrors.ValidationError{
				Validation: validation,
			},
		}

		if err != nil {
			response.Errors.Core = err.Error()
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	responseSuccess := func(user storage.User) error {
		s.authLog.Info(logger.USERCREATE, *params.Email, logger.USERSUCCESS)

		return ctx.JSON(http.StatusOK, api.User200{
			Success: api.Success{
				Base: api.Base{
					Version: s.Version,
				},
				Message: SuccessMessage,
			},
			Data: &[]api.UserParamsWithId{{
				ID:         user.ID,
				UserParams: userToUserParams(user),
			}},
		})
	}

	validation := map[string]interface{}{}

	if params.Email == nil && params.Mobile == nil {
		validation["email"] = "required"
		validation["mobile"] = "required"
	}

	if len(validation) > 0 {
		return responseBadRequest(nil, validation)
	}

	id := uuid.New().String()

	in := storage.User{
		ID: &id,
	}

	if params.Mobile != nil {
		in.Mobile = params.Mobile
	}

	if params.Email != nil {
		in.Email = params.Email
	}

	if params.FirstName != nil {
		in.FirstName = params.FirstName
	}

	if params.LastName != nil {
		in.LastName = params.LastName
	}

	if params.MiddleName != nil {
		in.MiddleName = params.MiddleName
	}

	u, err := s.storage.StoreUser(in)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrAlreadyExists):
			return responseBadRequest(err, map[string]interface{}{
				"mobile": "unique",
				"email":  "unique",
			})
		}

		return responseInternalServerError(err)
	}

	return responseSuccess(u)
}

// GET /users
func (s *Server) GetUsers(ctx echo.Context, params api.GetUsersParams) error {
	responseInternalServerError := func(err error) error {
		panic(err)
	}

	responseNotFound := func(err error) error {
		return s.error(ctx, http.StatusNotFound, err)
	}

	responseSuccess := func(users []storage.User) error {
		var userParamsWithIDs []api.UserParamsWithId
		for _, user := range users {
			userParamsWithIDs = append(userParamsWithIDs, api.UserParamsWithId{
				ID:         user.ID,
				UserParams: userToUserParams(user),
			})
		}

		return ctx.JSON(http.StatusOK, api.User200{
			Success: api.Success{
				Base: api.Base{
					Version: s.Version,
				},
				Errors:  nil,
				Message: SuccessMessage,
			},
			Data: &userParamsWithIDs,
		})
	}

	var users []storage.User
	var err error

	if params.Phone != nil || params.Email != nil {
		var user storage.User
		user, err = s.storage.GetUser(storage.User{
			Email:  (*string)(params.Email),
			Mobile: (*string)(params.Phone),
		})
		users = []storage.User{user}
	} else {
		users, err = s.storage.GetUsers()
	}

	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return responseNotFound(err)
		}

		return responseInternalServerError(err)
	}

	return responseSuccess(users)
}

// GET /users/{userID}
func (s *Server) GetUser(ctx echo.Context, userID api.UserID) error {
	responseInternalServerError := func(err error) error {
		panic(err)
	}

	responseNotFound := func(err error) error {
		return s.error(ctx, http.StatusNotFound, err)
	}

	responseBadRequest := func(err error, validation map[string]interface{}) error {
		response := apiErrors.ValidationErrorResponse{
			Response: apiErrors.Response{
				Version: s.Version,
				Message: apiErrors.ValidationErrorMessage,
			},
			Errors: &apiErrors.ValidationError{
				Validation: validation,
			},
		}

		if err != nil {
			response.Errors.Core = err.Error()
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	responseSuccess := func(user storage.User) error {
		return ctx.JSON(http.StatusOK, api.User200{
			Success: api.Success{
				Base: api.Base{
					Version: s.Version,
				},
				Message: SuccessMessage,
			},
			Data: &[]api.UserParamsWithId{{
				ID:         user.ID,
				UserParams: userToUserParams(user),
			}},
		})
	}

	if !strfmt.IsUUID(string(userID)) {
		return responseBadRequest(nil, map[string]interface{}{
			"userID": "format",
		})
	}

	id := string(userID)
	u, err := s.storage.GetUser(storage.User{
		ID: &id,
	})
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return responseNotFound(err)
		}

		return responseInternalServerError(err)
	}

	return responseSuccess(u)
}

// PATCH /users/{userID}
func (s *Server) PatchUser(ctx echo.Context, userID api.UserID) error {
	var params api.PatchUserJSONRequestBody
	if err := ctx.Bind(&params); err != nil {
		return s.error(ctx, http.StatusBadRequest, err)
	}

	responseInternalServerError := func(err error) error {
		s.authLog.Info(logger.USERUPDATE, *params.Email, logger.USERFAIL)

		panic(err)
	}

	responseNotFound := func(err error) error {
		s.authLog.Info(logger.USERUPDATE, *params.Email, logger.USERFAIL)

		return s.error(ctx, http.StatusNotFound, err)
	}

	responseBadRequest := func(err error, validation map[string]interface{}) error {
		s.authLog.Info(logger.USERUPDATE, *params.Email, logger.USERFAIL)

		response := apiErrors.ValidationErrorResponse{
			Response: apiErrors.Response{
				Version: s.Version,
				Message: apiErrors.ValidationErrorMessage,
			},
			Errors: &apiErrors.ValidationError{
				Validation: validation,
			},
		}

		if err != nil {
			response.Errors.Core = err.Error()
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	responseSuccess := func(user storage.User) error {
		s.authLog.Info(logger.USERUPDATE, *params.Email, logger.USERSUCCESS)

		return ctx.JSON(http.StatusOK, api.User200{
			Success: api.Success{
				Base: api.Base{
					Version: s.Version,
				},
				Message: SuccessMessage,
			},
			Data: &[]api.UserParamsWithId{{
				ID:         user.ID,
				UserParams: userToUserParams(user),
			}},
		})
	}

	if !strfmt.IsUUID(string(userID)) {
		return responseBadRequest(nil, map[string]interface{}{
			"userID": "format",
		})
	}

	in := storage.User{
		FirstName:  params.FirstName,
		LastName:   params.LastName,
		MiddleName: params.MiddleName,
		Mobile:     params.Mobile,
		Email:      params.Email,
	}

	u, err := s.storage.UpdateUser(string(userID), in)
	if err != nil {
		switch err {
		case storage.ErrUserAlreadyExists:
			return responseBadRequest(err, nil)
		case storage.ErrUserNotFound:
			return responseNotFound(err)
		case storage.ErrUserContactBelongsToOtherUser:
			return responseBadRequest(err, nil)
		default:
			return responseInternalServerError(err)
		}
	}

	return responseSuccess(u)
}

// DELETE /users/{userID}
func (s *Server) DeleteUser(ctx echo.Context, userID api.UserID) error {
	responseInternalServerError := func(err error) error {
		s.authLog.Info(logger.USERDELETE, string(userID), logger.USERFAIL)

		panic(err)
	}

	responseNotFound := func(err error) error {
		s.authLog.Info(logger.USERDELETE, string(userID), logger.USERFAIL)

		return s.error(ctx, http.StatusNotFound, err)
	}

	responseBadRequest := func(err error, validation map[string]interface{}) error {
		s.authLog.Info(logger.USERDELETE, string(userID), logger.USERFAIL)

		response := apiErrors.ValidationErrorResponse{
			Response: apiErrors.Response{
				Version: s.Version,
				Message: apiErrors.ValidationErrorMessage,
			},
			Errors: &apiErrors.ValidationError{
				Validation: validation,
			},
		}

		if err != nil {
			response.Errors.Core = err.Error()
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	responseSuccess := func() error {
		s.authLog.Info(logger.USERDELETE, string(userID), logger.USERSUCCESS)

		return ctx.JSON(http.StatusOK, api.User200{
			Success: api.Success{
				Base: api.Base{
					Version: s.Version,
				},
				Message: SuccessMessage,
			},
			Data: &[]api.UserParamsWithId{},
		})
	}

	if !strfmt.IsUUID(string(userID)) {
		return responseBadRequest(nil, map[string]interface{}{
			"userID": "format",
		})
	}

	err := s.storage.DeleteUser(string(userID))
	if err != nil {
		switch err {
		case storage.ErrUserNotFound:
			return responseNotFound(err)
		default:
			return responseInternalServerError(err)
		}
	}

	return responseSuccess()
}

func userToUserParams(user storage.User) api.UserParams {
	params := api.UserParams{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Email:      user.Email,
		Mobile:     user.Mobile,
	}

	return params
}
