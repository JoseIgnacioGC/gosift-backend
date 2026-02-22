package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func formatValidationErrors(err error) []string {
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

func register(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequestDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": formatValidationErrors(err)})
			return
		}

		resp, err := service.Register(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, ErrEmailAlreadyExists) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}
