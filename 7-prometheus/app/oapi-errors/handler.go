package apierrors

import (
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"

	middleware "github.com/foxcool/homework/7-prometheus/app/oapi-validator"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	version string

	sentryHub *sentry.Hub
}

func NewHandler(version string) *Handler {
	return &Handler{
		version: version,
	}
}

func (h *Handler) WithSentry(hub *sentry.Hub) *Handler {
	return &Handler{
		version: h.version,

		sentryHub: hub,
	}
}

func (h *Handler) Handle(err error, ctx echo.Context) {
	code := http.StatusInternalServerError

	if requestError, ok := err.(*openapi3filter.RequestError); ok {
		code = requestError.HTTPStatus()

		// Get original error
		err = requestError.Err
	}

	switch err := err.(type) {
	case *echo.HTTPError:
		code = err.Code
	case *openapi3filter.RouteError:
		if err.Error() == "Path was not found" {
			response := NotFoundResponse{
				Response: Response{
					Version: h.version,
					Message: UnknownPathMessage,
				},
				Errors: []interface{}{
					err.Error(),
				},
			}

			_ = ctx.JSON(http.StatusNotFound, response)
		}

		if err.Error() == "Path doesn't support the HTTP method" {
			response := MethodNotAllowedResponse{
				Response: Response{
					Version: h.version,
					Message: fmt.Sprintf(MethodNotAllowedMessage, ctx.Request().Method),
				},
				Errors: []interface{}{
					err.Error(),
				},
			}

			_ = ctx.JSON(http.StatusMethodNotAllowed, response)
		}
		return
	case *middleware.CompositeError:
		if err.Type != middleware.CompositeErrorTypeValidation {
			break
		}

		var core []string
		validation := make(map[string]interface{})
		for _, err := range err.Errors() {
			switch err := err.(type) {
			case *middleware.SchemaError:
				if err.Origin.Reason != "" {
					core = append(core, err.Origin.Reason)
				}

				convert := func(schemaField string) string {
					switch err.Origin.SchemaField {
					case "minLength":
						return "min"
					case "maxLength":
						return "max"
					case "type":
						return err.Origin.Schema.Type
					}

					return schemaField
				}

				if strings.Contains(err.Key, ".") {
					keyParts := strings.Split(err.Key, ".")

					for i, keyPart := range keyParts {
						var validationPart map[string]interface{}

						for j := 0; j < i; j++ {
							if validation[keyParts[j]] == nil {
								validation[keyParts[j]] = map[string]interface{}{}
							}

							validationPart = validation[keyParts[j]].(map[string]interface{})
						}

						if i+1 == len(keyParts) && validationPart != nil {
							validationPart[keyPart] = convert(err.Origin.SchemaField)
						}
					}
				} else {
					validation[err.Key] = convert(err.Origin.SchemaField)
				}
			}
		}

		response := ValidationErrorResponse{
			Response: Response{
				Version: h.version,
				Message: ValidationErrorMessage,
			},
			Errors: &ValidationError{
				Core:       strings.Join(core, "\n"),
				Validation: validation,
			},
		}

		_ = ctx.JSON(http.StatusBadRequest, response)
		return
	case *openapi3.SchemaError:
		validationError := parseValidationErrors(err.Error())

		response := ValidationErrorResponse{
			Response: Response{
				Version: h.version,
				Message: ValidationErrorMessage,
			},
			Errors: validationError,
		}

		_ = ctx.JSON(http.StatusBadRequest, response)
		return
	}

	switch code {
	case http.StatusBadRequest:
		response := ValidationErrorResponse{
			Response: Response{
				Version: h.version,
				Message: ValidationErrorMessage,
			},
			Errors: &ValidationError{},
		}

		if err.Error() == "EOF" {
			response.Errors.JSON = "EOF"
		} else {
			response.Errors.Core = err.Error()
		}

		_ = ctx.JSON(http.StatusBadRequest, response)
	case http.StatusInternalServerError:
		h.capture(ctx.Request(), err)

		response := InternalServerErrorResponse{
			Response: Response{
				Version: h.version,
				Message: InternalServerErrorMessage,
			},
		}

		if err != nil {
			response.Errors = []interface{}{
				err.Error(),
			}
		}

		_ = ctx.JSON(http.StatusInternalServerError, response)
	case http.StatusNotFound:
		response := NotFoundResponse{
			Response: Response{
				Version: h.version,
				Message: NotFoundMessage,
			},
		}

		_ = ctx.JSON(http.StatusNotFound, response)
	case http.StatusMethodNotAllowed:
		response := MethodNotAllowedResponse{
			Response: Response{
				Version: h.version,
				Message: fmt.Sprintf(MethodNotAllowedMessage, ctx.Request().Method),
			},
			Errors: []interface{}{
				err.Error(),
			},
		}

		_ = ctx.JSON(http.StatusMethodNotAllowed, response)
	}

	// fmt.Printf("Error: (%T) %v \n", err, err)
}

func (h *Handler) capture(r *http.Request, err error) {
	if h.sentryHub == nil {
		return
	}

	stacktrace := sentry.NewStacktrace()

	if len(stacktrace.Frames) > 4 {
		stacktrace.Frames = stacktrace.Frames[:len(stacktrace.Frames)-4]
	}

	event := sentry.Event{
		User:    sentry.User{},
		Request: sentryRequestFromHTTP(r),
		Exception: []sentry.Exception{{
			Type:       fmt.Sprintf("%T", err),
			Value:      err.Error(),
			Stacktrace: stacktrace,
		}},
	}

	eventID := h.sentryHub.CaptureEvent(&event)

	log.WithError(err).
		WithFields(log.Fields{
			"sentryEventID": eventID,
		}).
		Error("Panic captured!")

	if log.StandardLogger().GetLevel() == log.DebugLevel {
		debug.PrintStack()
	}
}

func sentryRequestFromHTTP(r *http.Request) *sentry.Request {
	proto := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		proto = "https"
	}

	sentryRequest := sentry.Request{
		URL:         proto + "://" + r.Host + r.URL.Path,
		Method:      r.Method,
		QueryString: r.RequestURI,
		Cookies:     r.Header.Get("Cookie"),
		Headers:     map[string]string{},
	}

	for k, v := range r.Header {
		sentryRequest.Headers[k] = strings.Join(v, ",")
	}

	if addr, port, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		sentryRequest.Env = map[string]string{
			"REMOTE_ADDR": addr,
			"REMOTE_PORT": port,
		}
	}

	return &sentryRequest
}
