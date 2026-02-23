package auth

import (
	"errors"
	"net/http"

	"github.com/JoseIgnacioGC/gosift-backend/internal/validation"
	"github.com/gin-gonic/gin"
)

func register(service *service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequestDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": validation.FormatErrors(err)})
			return
		}

		resp, err := service.register(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, errEmailAlreadyExists) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

func login(service *service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequestDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": validation.FormatErrors(err)})
			return
		}

		resp, err := service.login(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
