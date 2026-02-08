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

// FileMetadata represents file metadata
type FileMetadata struct {
	ID          string `json:"id" example:"507f1f77bcf86cd799439011"`
	Filename    string `json:"filename" example:"document.pdf"`
	Size        int64  `json:"size" example:"1024000"`
	ContentType string `json:"content_type" example:"application/pdf"`
	CreatedAt   int64  `json:"created_at" example:"1704067200"`
}

// FileResponse represents the response for file upload
type FileResponse struct {
	File FileMetadata `json:"file"`
}

// ListFilesResponse represents the response for listing files
type ListFilesResponse struct {
	Files []FileMetadata `json:"files"`
}

// GetFileResponse represents the response for getting a single file
type GetFileResponse struct {
	File FileMetadata `json:"file"`
}

// DeleteFileResponse represents the response for deleting a file
type DeleteFileResponse struct {
	Success bool `json:"success" example:"true"`
}

// InitiateMultipartUploadRequest represents the request to initiate a multipart upload
type InitiateMultipartUploadRequest struct {
	Filename    string `json:"filename" binding:"required" example:"large-video.mp4"`
	ContentType string `json:"content_type" example:"video/mp4"`
	TotalSize   int64  `json:"total_size" binding:"required" example:"1073741824"`
}

// InitiateMultipartUploadResponse represents the response for initiating a multipart upload
type InitiateMultipartUploadResponse struct {
	UploadID   string `json:"upload_id" example:"abc123"`
	ChunkSize  int64  `json:"chunk_size" example:"10485760"`
	TotalParts int32  `json:"total_parts" example:"10"`
}

// UploadPartResponse represents the response for uploading a single part
type UploadPartResponse struct {
	Etag       string `json:"etag" example:"\"etag1\""`
	PartNumber int32  `json:"part_number" example:"1"`
}

// PartInfo represents a single part in the complete request
type PartInfo struct {
	PartNumber int32  `json:"part_number" example:"1"`
	Etag       string `json:"etag" example:"\"etag1\""`
}

// CompleteMultipartUploadRequest represents the request to complete a multipart upload
type CompleteMultipartUploadRequest struct {
	Parts []PartInfo `json:"parts" binding:"required"`
}

// AbortMultipartUploadResponse represents the response for aborting a multipart upload
type AbortMultipartUploadResponse struct {
	Success bool `json:"success" example:"true"`
}
