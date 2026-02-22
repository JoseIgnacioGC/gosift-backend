package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatErrors(err error) []string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return []string{err.Error()}
	}

	errs := make([]string, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		errs = append(errs, fmt.Sprintf("%s failed on '%s' validation, got '%v'", fieldErr.Field(), fieldErr.Tag(), fieldErr.Value()))
	}

	return errs
}
