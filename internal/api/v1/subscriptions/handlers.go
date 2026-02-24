package subscriptions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JoseIgnacioGC/gosift-backend/internal/validation"
	"github.com/JoseIgnacioGC/gosift-backend/internal/middleware"
)

func createSubscription(svc *service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.UserIDKey)

		var req CreateRequestDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": validation.FormatErrors(err)})
			return
		}

		resp, err := svc.create(c.Request.Context(), userID, req)
		if err != nil {
			if errors.Is(err, errFeedAlreadySubscribed) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

func listSubscriptions(svc *service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.UserIDKey)

		subs, err := svc.list(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
			return
		}

		c.JSON(http.StatusOK, subs)
	}
}

func updateSubscription(svc *service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.UserIDKey)
		subID := c.Param("id")

		var req UpdateRequestDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": validation.FormatErrors(err)})
			return
		}

		if err := svc.update(c.Request.Context(), userID, subID, req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

func deleteSubscription(svc *service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.UserIDKey)
		subID := c.Param("id")

		if err := svc.delete(c.Request.Context(), userID, subID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}
