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
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx := h.contextWithAuth(c)
	resp, err := h.client.ListFiles(ctx, &filev1.ListFilesRequest{
		UserId: user.Id,
	})
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

func (h *FileHandler) GetFile(c *gin.Context) {
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
	resp, err := h.client.GetFile(ctx, &filev1.GetFileRequest{
		Id:     fileID,
		UserId: user.Id,
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
			return
		}

		chunk := msg.GetChunk()
		if chunk != nil {
			c.Writer.Write(chunk)
			c.Writer.Flush()
		}
	}
}
