package handlers

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"testing"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// SignUpRequest represents the signup request body
type SignUpRequest struct {
	Username string `json:"username" binding:"required" example:"testing"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// UserResponse represents user information in responses
type UserResponse struct {
	ID       string `json:"id" example:"69654eb7a1135a809430d0b7"`
	Username string `json:"username" example:"testing"`
	Role     string `json:"role" example:"ROLE_USER"`
}

// AuthResponse represents the response for login and signup
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"jwt-token-123"`
}

// DeleteUserRequest represents the delete user request body
type DeleteUserRequest struct {
	ID string `json:"id" binding:"required" example:"69654eb7a1135a809430d0b7"`
}

// DeleteUserResponse represents the delete user response
type DeleteUserResponse struct {
	Success bool `json:"success" example:"true"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}
