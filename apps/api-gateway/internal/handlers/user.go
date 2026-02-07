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

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Admin-only endpoint to delete a user account. Cannot delete own account or other admin accounts.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request body DeleteUserRequest true "User ID to delete"
// @Success      200 {object} DeleteUserResponse "User deleted successfully"
// @Failure      400 {object} ErrorResponse "Invalid request body or invalid user ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      403 {object} ErrorResponse "Forbidden - cannot delete own account or admin accounts"
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/admin/delete_user [delete]
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

// ListUsers godoc
// @Summary      List all users
// @Description  Admin-only endpoint to list all users with optional role filter and username search
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        role query string false "Filter by role (admin or user)"
// @Param        username query string false "Search by username (partial match)"
// @Success      200 {object} map[string][]map[string]interface{} "List of users"
// @Failure      400 {object} ErrorResponse "Invalid role parameter"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/admin/list_users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	roleParam := c.Query("role")
	username := c.Query("username")

	var role userv1.Role
	if roleParam != "" {
		switch roleParam {
		case "admin":
			role = userv1.Role_ROLE_ADMIN
		case "user":
			role = userv1.Role_ROLE_USER
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "role must be 'admin' or 'user'"})
			return
		}
	} else {
		role = userv1.Role_ROLE_UNSPECIFIED
	}

	resp, err := h.client.ListUsers(c, &userv1.ListUsersRequest{
		Role:           role,
		UsernameFilter: username,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
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

	users := make([]map[string]interface{}, len(resp.Users))
	for i, user := range resp.Users {
		users[i] = map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"role":     user.Role.String(),
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
