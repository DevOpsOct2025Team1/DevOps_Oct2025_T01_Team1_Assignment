package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

type UserHandler struct {
	client UserServiceClient
}

func NewUserHandler(client UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	var req struct {
		Id string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.DeleteAccount(c, &userv1.DeleteUserByIdRequest{Id: req.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": resp.Success,
	})
}
