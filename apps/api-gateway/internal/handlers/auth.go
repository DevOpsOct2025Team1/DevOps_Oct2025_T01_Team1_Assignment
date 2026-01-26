package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
)

type AuthHandler struct {
	client AuthServiceClient
}

func NewAuthHandler(client AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

// SignUp godoc
// @Summary      Create a new user
// @Description  Admin-only endpoint to create a new user account
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request body SignUpRequest true "User credentials"
// @Success      200 {object} AuthResponse "User created successfully"
// @Failure      400 {object} ErrorResponse "Invalid request body"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      403 {object} ErrorResponse "Forbidden - admin role required"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/admin/create_user [post]
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and receive a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "User credentials"
// @Success      200 {object} AuthResponse "Login successful"
// @Failure      400 {object} ErrorResponse "Invalid request body"
// @Failure      500 {object} ErrorResponse "Invalid credentials or server error"
// @Router       /api/login [post]
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
