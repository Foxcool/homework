package apierrors

import (
	"regexp"
)

var (
	patternRequired = regexp.MustCompile(`Property '(\w+)' is missing`)
	patternFormat   = regexp.MustCompile(`Field must be set to (\w+) or not be present`)
)

func parseValidationErrors(err string) *ValidationError {
	res := &ValidationError{
		Validation: map[string]interface{}{},
	}

	wrapCore := func(message string) {
		if res.Core != "" {
			res.Core += "\n"
		}

		res.Core += message
	}

	if m := patternRequired.FindStringSubmatch(err); m != nil {
		res.Validation[m[1]] = "required"
	} else if m := patternFormat.FindStringSubmatch(err); m != nil {
		res.Validation[m[1]] = "format"
	} else {
		wrapCore(err)
	}

	return res
}
