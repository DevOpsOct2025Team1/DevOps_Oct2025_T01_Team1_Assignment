package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthHandler struct {
	authServiceAddr string
}

func NewAuthHandler(authServiceAddr string) *AuthHandler {
	return &AuthHandler{
		authServiceAddr: authServiceAddr,
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

	conn, err := grpc.NewClient(h.authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to auth service"})
		return
	}
	defer conn.Close()

	client := authv1.NewAuthServiceClient(conn)
	resp, err := client.SignUp(context.Background(), &authv1.SignUpRequest{
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := grpc.NewClient(h.authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to auth service"})
		return
	}
	defer conn.Close()

	client := authv1.NewAuthServiceClient(conn)
	resp, err := client.Login(context.Background(), &authv1.LoginRequest{
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

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := grpc.NewClient(h.authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to auth service"})
		return
	}
	defer conn.Close()

	client := authv1.NewAuthServiceClient(conn)
	resp, err := client.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: req.Token,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user": gin.H{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"role":     resp.User.Role,
		},
	})
}
