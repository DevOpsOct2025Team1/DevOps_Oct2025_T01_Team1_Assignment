package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	client AuthServiceClient
}

func NewAuthHandler(client AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.SignUp(c, &authv1.SignUpRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				statusCode = http.StatusBadRequest
			case codes.Unauthenticated:
				statusCode = http.StatusUnauthorized
			case codes.NotFound:
				statusCode = http.StatusNotFound
			case codes.AlreadyExists:
				statusCode = http.StatusConflict
			case codes.PermissionDenied:
				statusCode = http.StatusForbidden
			}
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"role":     resp.User.Role,
		},
		"token": resp.Token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.Login(c, &authv1.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				statusCode = http.StatusBadRequest
			case codes.Unauthenticated:
				statusCode = http.StatusUnauthorized
			case codes.NotFound:
				statusCode = http.StatusNotFound
			case codes.AlreadyExists:
				statusCode = http.StatusConflict
			case codes.PermissionDenied:
				statusCode = http.StatusForbidden
			}
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"role":     resp.User.Role,
		},
		"token": resp.Token,
	})
}
