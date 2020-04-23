package apierrors

var (
	InternalServerErrorMessage = "Internal server error"
	ValidationErrorMessage     = "Validation error"
	NotFoundMessage            = "Entity not found"
	AccessDeniedMessage        = "Access denied"
	UnknownPathMessage         = "Unknown path"
	MethodNotAllowedMessage    = "Method %s not allowed"
)

type Response struct {
	Version string `json:"version"`

	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

// Validation error (400)
type ValidationError struct {
	Core string `json:"core,omitempty"`
	JSON string `json:"json,omitempty"`

	Validation map[string]interface{} `json:"validation"`
}

type ValidationErrorResponse struct {
	Response

	Errors *ValidationError `json:"errors,omitempty"`
}

// Not found (404)
type NotFoundResponse struct {
	Response

	Errors []interface{} `json:"errors"`
}

// Method not allowed (405)
type MethodNotAllowedResponse struct {
	Response

	Errors []interface{} `json:"errors,omitempty"`
}

// Internal server error (500)
type InternalServerErrorResponse struct {
	Response

	Errors []interface{} `json:"errors"`
}
