package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/handlers"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

// ValidateRole validates the user's role before allowing requests to pass to other handlers
func ValidateRole(authService handlers.AuthServiceClient, roles []userv1.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		auth := strings.Split(authorization, " ")
		if len(auth) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}
		if auth[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization type"})
			c.Abort()
			return
		}

		resp, err := authService.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
			Token: auth[1],
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		role := resp.GetUser().GetRole()

		for _, requiredRole := range roles {
			if role == requiredRole {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
	}
}
