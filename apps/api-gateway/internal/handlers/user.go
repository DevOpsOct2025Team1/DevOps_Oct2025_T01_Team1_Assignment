package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	client UserServiceClient
}

func NewUserHandler(client UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var req struct {
		Id string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	currentUser := currentUserVal.(*userv1.User)

	if currentUser.Id == req.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete your own account"})
		return
	}

	targetUserResp, err := h.client.GetUser(c, &userv1.GetUserRequest{Id: req.Id})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if targetUserResp.User.Role == userv1.Role_ROLE_ADMIN {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete admin accounts"})
		return
	}

	resp, err := h.client.DeleteAccount(c, &userv1.DeleteUserByIdRequest{Id: req.Id})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": resp.Success,
	})
}
