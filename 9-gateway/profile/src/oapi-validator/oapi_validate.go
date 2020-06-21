package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

const (
	EchoContextKey = "oapi-codegen/echo-context"
	UserDataKey    = "oapi-codegen/user-data"
)

var (
	errRouteMissingSwagger = errors.New("route is missing OpenAPI specification")
	MetricPath             = "/metrics"
)

func OAPIRequestValidator(swagger *openapi3.Swagger) echo.MiddlewareFunc {
	return OAPIRequestValidatorWithOptions(swagger, nil)
}

type Options struct {
	Options      openapi3filter.Options
	ParamDecoder openapi3filter.ContentParameterDecoder
	UserData     interface{}
}

func OAPIRequestValidatorWithOptions(swagger *openapi3.Swagger, options *Options) echo.MiddlewareFunc {
	router := openapi3filter.NewRouter().WithSwagger(swagger)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := ValidateRequestFromContext(c, router, options)
			if err != nil {
				return err
			}

			return next(c)
		}
	}
}

func ValidateRequestFromContext(ctx echo.Context, router *openapi3filter.Router, options *Options) error {
	r := ctx.Request()

	if r.RequestURI == MetricPath {
		return nil
	}

	route, pathParams, err := router.FindRoute(r.Method, r.URL)

	// We failed to find a matching route for the request.
	if err != nil {
		switch e := err.(type) {
		case *openapi3filter.RouteError:
			// We've got a bad request, the path requested doesn't match
			// either server, or path, or something.
			return echo.NewHTTPError(http.StatusBadRequest, e.Reason)
		default:
			// This should never happen today, but if our upstream code changes,
			// we don't want to crash the server, so handle the unexpected error.
			return echo.NewHTTPError(http.StatusInternalServerError,
				fmt.Sprintf("error validating route: %s", err.Error()))
		}
	}

	validationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}

	requestContext := context.WithValue(context.Background(), EchoContextKey, ctx)
	if options != nil {
		validationInput.Options = &options.Options
		validationInput.ParamDecoder = options.ParamDecoder
		requestContext = context.WithValue(requestContext, UserDataKey, options.UserData)
	} else {
		validationInput.Options = openapi3filter.DefaultOptions
	}

	// Validate request
	operationParameters := validationInput.Route.Operation.Parameters
	pathItemParameters := route.PathItem.Parameters

	for _, parameterRef := range pathItemParameters {
		parameter := parameterRef.Value
		if operationParameters != nil {
			if override := operationParameters.GetByInAndName(parameter.In, parameter.Name); override != nil {
				continue
			}
			if err := openapi3filter.ValidateParameter(requestContext, validationInput, parameter); err != nil {
				return err
			}
		}
	}

	for _, parameter := range operationParameters {
		if err := openapi3filter.ValidateParameter(requestContext, validationInput, parameter.Value); err != nil {
			return err
		}
	}

	// Validate request body
	requestBody := validationInput.Route.Operation.RequestBody
	if requestBody != nil && !validationInput.Options.ExcludeRequestBody {
		if err := openapi3filter.ValidateRequestBody(requestContext, validationInput, requestBody.Value); err != nil {
			return err
		}
	}

	// Security
	security := validationInput.Route.Operation.Security

	if security == nil {
		if route.Swagger == nil {
			return errRouteMissingSwagger
		} else {
			security = &route.Swagger.Security
		}
	}

	if security != nil {
		if err := openapi3filter.ValidateSecurityRequirements(requestContext, validationInput, *security); err != nil {
			return err
		}
	}

	return nil
}
