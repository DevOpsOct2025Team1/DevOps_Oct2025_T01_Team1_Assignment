package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	filev1 "github.com/provsalt/DOP_P01_Team1/common/file/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type FileHandler struct {
	client FileServiceClient
}

func NewFileHandler(client FileServiceClient) *FileHandler {
	return &FileHandler{client: client}
}

func (h *FileHandler) contextWithAuth(c *gin.Context) context.Context {
	authHeader := c.GetHeader("Authorization")
	ctx := c.Request.Context()

	if authHeader != "" {
		md := metadata.Pairs("authorization", authHeader)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return ctx
}

func (h *FileHandler) getUserFromContext(c *gin.Context) (*userv1.User, error) {
	userVal, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not authenticated")
	}
	user, ok := userVal.(*userv1.User)
	if !ok {
		return nil, fmt.Errorf("invalid user context")
	}
	return user, nil
}

func mapGRPCError(err error) int {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			return http.StatusBadRequest
		case codes.NotFound:
			return http.StatusNotFound
		case codes.PermissionDenied:
			return http.StatusForbidden
		case codes.Unauthenticated:
			return http.StatusUnauthorized
		case codes.AlreadyExists:
			return http.StatusConflict
		case codes.ResourceExhausted:
			return http.StatusTooManyRequests
		case codes.Unavailable:
			return http.StatusServiceUnavailable
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

// ListFiles godoc
// @Summary      List all files for the authenticated user
// @Description  Retrieve a list of all files uploaded by the authenticated user
// @Tags         files
// @Produce      json
// @Success      200 {object} ListFilesResponse "Files retrieved successfully"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx := h.contextWithAuth(c)
	resp, err := h.client.ListFiles(ctx, &filev1.ListFilesRequest{})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	files := make([]map[string]interface{}, len(resp.Files))
	for i, file := range resp.Files {
		files[i] = map[string]interface{}{
			"id":           file.Id,
			"filename":     file.Filename,
			"size":         file.Size,
			"content_type": file.ContentType,
			"created_at":   file.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// GetFile godoc
// @Summary      Get file metadata
// @Description  Retrieve metadata for a specific file by ID
// @Tags         files
// @Produce      json
// @Param        id path string true "File ID"
// @Success      200 {object} GetFileResponse "File metadata retrieved successfully"
// @Failure      400 {object} ErrorResponse "Invalid file ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} ErrorResponse "File not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file id is required"})
		return
	}

	ctx := h.contextWithAuth(c)
	resp, err := h.client.GetFile(ctx, &filev1.GetFileRequest{
		Id: fileID,
	})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file": map[string]interface{}{
			"id":           resp.File.Id,
			"filename":     resp.File.Filename,
			"size":         resp.File.Size,
			"content_type": resp.File.ContentType,
			"created_at":   resp.File.CreatedAt,
		},
	})
}

// DeleteFile godoc
// @Summary      Delete a file
// @Description  Delete a file by ID (must be owned by the authenticated user)
// @Tags         files
// @Produce      json
// @Param        id path string true "File ID"
// @Success      200 {object} DeleteFileResponse "File deleted successfully"
// @Failure      400 {object} ErrorResponse "Invalid file ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} ErrorResponse "File not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file id is required"})
		return
	}

	ctx := h.contextWithAuth(c)
	resp, err := h.client.DeleteFile(ctx, &filev1.DeleteFileRequest{
		Id:     fileID,
		UserId: user.Id,
	})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success})
}

// UploadFile godoc
// @Summary      Upload a file
// @Description  Upload a file to S3 and save metadata for the authenticated user
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File to upload"
// @Success      200 {object} FileResponse "File uploaded successfully"
// @Failure      400 {object} ErrorResponse "Invalid file or missing required field"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      429 {object} ErrorResponse "Too many files - limit is 20 per user"
// @Failure      500 {object} ErrorResponse "Internal server error or file too large"
// @Security     BearerAuth
// @Router       /api/files [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	ctx := h.contextWithAuth(c)
	stream, err := h.client.UploadFile(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = stream.Send(&filev1.UploadFileRequest{
		Data: &filev1.UploadFileRequest_Metadata{
			Metadata: &filev1.UploadFileMetadata{
				Filename:    file.Filename,
				ContentType: file.Header.Get("Content-Type"),
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send metadata"})
		return
	}

	buffer := make([]byte, 64*1024)
	for {
		n, err := src.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		err = stream.Send(&filev1.UploadFileRequest{
			Data: &filev1.UploadFileRequest_Chunk{
				Chunk: buffer[:n],
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload chunk"})
			return
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file": map[string]interface{}{
			"id":           resp.File.Id,
			"filename":     resp.File.Filename,
			"size":         resp.File.Size,
			"content_type": resp.File.ContentType,
			"created_at":   resp.File.CreatedAt,
		},
	})
}

// DownloadFile godoc
// @Summary      Download a file
// @Description  Download a file by ID (must be owned by the authenticated user)
// @Tags         files
// @Produce      octet-stream
// @Param        id path string true "File ID"
// @Success      200 {file} binary "File content"
// @Failure      400 {object} ErrorResponse "Invalid file ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} ErrorResponse "File not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files/{id}/download [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file id is required"})
		return
	}

	ctx := h.contextWithAuth(c)
	stream, err := h.client.DownloadFile(ctx, &filev1.DownloadFileRequest{
		Id: fileID,
	})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	firstMsg, err := stream.Recv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to receive metadata"})
		return
	}

	metadata := firstMsg.GetMetadata()
	if metadata == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response: expected metadata"})
		return
	}

	c.Header("Content-Type", metadata.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", metadata.Filename))
	c.Header("Content-Length", strconv.FormatInt(metadata.Size, 10))

	c.Status(http.StatusOK)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Headers already sent, log the error and abort the connection
			fmt.Printf("Error streaming file download: %v\n", err)
			return
		}

		chunk := msg.GetChunk()
		if chunk != nil {
			c.Writer.Write(chunk)
			c.Writer.Flush()
		}
	}
}

// InitiateMultipartUpload godoc
// @Summary      Initiate a multipart upload
// @Description  Start a new multipart upload session for large files
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        body body InitiateMultipartUploadRequest true "Upload initiation request"
// @Success      200 {object} InitiateMultipartUploadResponse "Upload session created"
// @Failure      400 {object} ErrorResponse "Invalid request body"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files/multipart/initiate [post]
func (h *FileHandler) InitiateMultipartUpload(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type"`
		TotalSize   int64  `json:"total_size" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx := h.contextWithAuth(c)
	resp, err := h.client.InitiateMultipartUpload(ctx, &filev1.InitiateMultipartUploadRequest{
		Filename:    req.Filename,
		ContentType: req.ContentType,
		TotalSize:   req.TotalSize,
	})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_id":   resp.UploadId,
		"chunk_size":  resp.ChunkSize,
		"total_parts": resp.TotalParts,
	})
}

// UploadPart godoc
// @Summary      Upload a part of a multipart upload
// @Description  Upload a single chunk of a file as part of a multipart upload
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        upload_id path string true "Upload session ID"
// @Param        part_number path int true "Part number"
// @Param        chunk formData file true "File chunk to upload"
// @Success      200 {object} UploadPartResponse "Part uploaded successfully"
// @Failure      400 {object} ErrorResponse "Invalid part number or missing chunk"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files/multipart/{upload_id}/part/{part_number} [post]
func (h *FileHandler) UploadPart(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	uploadID := c.Param("upload_id")
	partNumberStr := c.Param("part_number")
	partNumber, err := strconv.Atoi(partNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid part number"})
		return
	}

	file, err := c.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chunk is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open chunk"})
		return
	}
	defer src.Close()

	chunk, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read chunk"})
		return
	}

	ctx := h.contextWithAuth(c)
	resp, err := h.client.UploadPart(ctx, &filev1.UploadPartRequest{
		UploadId:   uploadID,
		PartNumber: int32(partNumber),
		Chunk:      chunk,
	})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"etag":        resp.Etag,
		"part_number": resp.PartNumber,
	})
}

// CompleteMultipartUpload godoc
// @Summary      Complete a multipart upload
// @Description  Finalize a multipart upload by providing all part ETags
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        upload_id path string true "Upload session ID"
// @Param        body body CompleteMultipartUploadRequest true "Parts to complete"
// @Success      200 {object} FileResponse "File uploaded successfully"
// @Failure      400 {object} ErrorResponse "Invalid request body"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} ErrorResponse "Upload session not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files/multipart/{upload_id}/complete [post]
func (h *FileHandler) CompleteMultipartUpload(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	uploadID := c.Param("upload_id")

	var req struct {
		Parts []struct {
			PartNumber int32  `json:"part_number"`
			Etag       string `json:"etag"`
		} `json:"parts" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	parts := make([]*filev1.PartInfo, len(req.Parts))
	for i, p := range req.Parts {
		parts[i] = &filev1.PartInfo{
			PartNumber: p.PartNumber,
			Etag:       p.Etag,
		}
	}

	ctx := h.contextWithAuth(c)
	resp, err := h.client.CompleteMultipartUpload(ctx, &filev1.CompleteMultipartUploadRequest{
		UploadId: uploadID,
		Parts:    parts,
	})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file": map[string]interface{}{
			"id":           resp.File.Id,
			"filename":     resp.File.Filename,
			"size":         resp.File.Size,
			"content_type": resp.File.ContentType,
			"created_at":   resp.File.CreatedAt,
		},
	})
}

// AbortMultipartUpload godoc
// @Summary      Abort a multipart upload
// @Description  Cancel an in-progress multipart upload and clean up resources
// @Tags         files
// @Produce      json
// @Param        upload_id path string true "Upload session ID"
// @Success      200 {object} AbortMultipartUploadResponse "Upload aborted successfully"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/files/multipart/{upload_id} [delete]
func (h *FileHandler) AbortMultipartUpload(c *gin.Context) {
	_, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	uploadID := c.Param("upload_id")

	ctx := h.contextWithAuth(c)
	resp, err := h.client.AbortMultipartUpload(ctx, &filev1.AbortMultipartUploadRequest{
		UploadId: uploadID,
	})
	if err != nil {
		c.JSON(mapGRPCError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success})
}
