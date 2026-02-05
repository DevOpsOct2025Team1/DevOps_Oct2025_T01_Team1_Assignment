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
// @Param        request body object true "Delete user request" example({"id":"user123"})
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
		Id string `json:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	// Log the ID being deleted for debugging
	c.Request.Header.Set("X-Delete-User-ID", req.Id)

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

	targetUserResp, err := h.client.GetUser(c.Request.Context(), &userv1.GetUserRequest{Id: req.Id})
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

	resp, err := h.client.DeleteAccount(c.Request.Context(), &userv1.DeleteUserByIdRequest{Id: req.Id})
	if err != nil {
		// Log the error for debugging
		c.Error(err)
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
// @Description  Admin-only endpoint to retrieve a list of all users in the system
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200 {object} ListUsersResponse "List of users retrieved successfully"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      403 {object} ErrorResponse "Forbidden - admin role required"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/admin [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	resp, err := h.client.ListUsers(c, &userv1.ListUsersRequest{})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				statusCode = http.StatusUnauthorized
			case codes.PermissionDenied:
				statusCode = http.StatusForbidden
			default:
				statusCode = http.StatusInternalServerError
			}
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	users := make([]gin.H, 0, len(resp.Users))
	for _, u := range resp.Users {
		role := "user"
		if u.Role == userv1.Role_ROLE_ADMIN {
			role = "admin"
		}
		users = append(users, gin.H{
			"id":       u.Id,
			"username": u.Username,
			"role":     role,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
