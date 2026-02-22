package auth

import (
	"errors"
	"net/http"

	"github.com/JoseIgnacioGC/gosift-backend/internal/validation"
	"github.com/gin-gonic/gin"
)

func register(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequestDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": validation.FormatErrors(err)})
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
