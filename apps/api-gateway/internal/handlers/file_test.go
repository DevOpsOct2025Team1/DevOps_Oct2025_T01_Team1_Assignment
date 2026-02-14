package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	filev1 "github.com/provsalt/DOP_P01_Team1/common/file/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func setupFileTestRouter(handler *FileHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("user", &userv1.User{
			Id:       "user-123",
			Username: "testuser",
			Role:     userv1.Role_ROLE_USER,
		})
		c.Next()
	})

	router.GET("/api/files", handler.ListFiles)
	router.GET("/api/files/:id", handler.GetFile)
	router.DELETE("/api/files/:id", handler.DeleteFile)
	router.POST("/api/files", handler.UploadFile)
	router.GET("/api/files/:id/download", handler.DownloadFile)
	router.POST("/api/files/multipart/initiate", handler.InitiateMultipartUpload)
	router.POST("/api/files/multipart/:upload_id/part/:part_number", handler.UploadPart)
	router.POST("/api/files/multipart/:upload_id/complete", handler.CompleteMultipartUpload)
	router.DELETE("/api/files/multipart/:upload_id", handler.AbortMultipartUpload)

	return router
}

func TestListFiles_Success(t *testing.T) {
	mockClient := &mockFileClient{
		listFilesFunc: func(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error) {
			return &filev1.ListFilesResponse{
				Files: []*filev1.File{
					{
						Id:          "file-1",
						Filename:    "test.txt",
						Size:        1024,
						ContentType: "text/plain",
						CreatedAt:   1704067200,
					},
				},
			}, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	files := response["files"].([]interface{})
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}

func TestListFiles_NoUser(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/files", handler.ListFiles)

	req, _ := http.NewRequest("GET", "/api/files", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestListFiles_ServiceError(t *testing.T) {
	mockClient := &mockFileClient{
		listFilesFunc: func(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error) {
			return nil, status.Error(codes.Internal, "db error")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestListFiles_PermissionDenied(t *testing.T) {
	mockClient := &mockFileClient{
		listFilesFunc: func(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error) {
			return nil, status.Error(codes.PermissionDenied, "not authorized")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetFile_Success(t *testing.T) {
	mockClient := &mockFileClient{
		getFileFunc: func(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
			return &filev1.FileResponse{
				File: &filev1.File{
					Id:          "file-1",
					Filename:    "test.txt",
					Size:        1024,
					ContentType: "text/plain",
					CreatedAt:   1704067200,
				},
			}, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/file-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetFile_NotFound(t *testing.T) {
	mockClient := &mockFileClient{
		getFileFunc: func(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
			return nil, status.Error(codes.NotFound, "file not found")
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/missing-file", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteFile_Success(t *testing.T) {
	mockClient := &mockFileClient{
		deleteFileFunc: func(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
			return &filev1.DeleteFileResponse{
				Success: true,
			}, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("DELETE", "/api/files/file-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Errorf("expected success true, got %v", response["success"])
	}
}

func TestUploadFile_Success(t *testing.T) {
	mockStream := &mockUploadStream{
		sendFunc: func(req *filev1.UploadFileRequest) error {
			return nil
		},
		closeAndRecvFunc: func() (*filev1.FileResponse, error) {
			return &filev1.FileResponse{
				File: &filev1.File{
					Id:          "file-123",
					Filename:    "upload.txt",
					Size:        13,
					ContentType: "text/plain",
					CreatedAt:   1704067200,
				},
			}, nil
		},
	}

	mockClient := &mockFileClient{
		uploadFileFunc: func(ctx context.Context) (filev1.FileService_UploadFileClient, error) {
			return mockStream, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "upload.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestUploadFile_MissingFile(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("POST", "/api/files", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

type mockUploadStream struct {
	sendFunc         func(*filev1.UploadFileRequest) error
	closeAndRecvFunc func() (*filev1.FileResponse, error)
}

func (m *mockUploadStream) Send(req *filev1.UploadFileRequest) error {
	if m.sendFunc != nil {
		return m.sendFunc(req)
	}
	return nil
}

func (m *mockUploadStream) CloseAndRecv() (*filev1.FileResponse, error) {
	if m.closeAndRecvFunc != nil {
		return m.closeAndRecvFunc()
	}
	return nil, nil
}

func (m *mockUploadStream) Header() (metadata.MD, error) { return nil, nil }
func (m *mockUploadStream) Trailer() metadata.MD         { return nil }
func (m *mockUploadStream) CloseSend() error             { return nil }
func (m *mockUploadStream) Context() context.Context     { return context.Background() }
func (m *mockUploadStream) SendMsg(interface{}) error    { return nil }
func (m *mockUploadStream) RecvMsg(interface{}) error    { return io.EOF }

type mockDownloadStream struct {
	recvFunc func() (*filev1.DownloadFileResponse, error)
}

func (m *mockDownloadStream) Recv() (*filev1.DownloadFileResponse, error) {
	if m.recvFunc != nil {
		return m.recvFunc()
	}
	return nil, io.EOF
}

func (m *mockDownloadStream) Header() (metadata.MD, error) { return nil, nil }
func (m *mockDownloadStream) Trailer() metadata.MD         { return nil }
func (m *mockDownloadStream) CloseSend() error             { return nil }
func (m *mockDownloadStream) Context() context.Context     { return context.Background() }
func (m *mockDownloadStream) SendMsg(interface{}) error    { return nil }
func (m *mockDownloadStream) RecvMsg(interface{}) error    { return io.EOF }

func TestInitiateMultipartUpload_Success(t *testing.T) {
	mockClient := &mockFileClient{
		initiateMultipartFunc: func(ctx context.Context, req *filev1.InitiateMultipartUploadRequest) (*filev1.InitiateMultipartUploadResponse, error) {
			if req.Filename != "big.mp4" {
				t.Errorf("expected filename big.mp4, got %s", req.Filename)
			}
			return &filev1.InitiateMultipartUploadResponse{
				UploadId:   "upload-abc",
				ChunkSize:  10 * 1024 * 1024,
				TotalParts: 5,
			}, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := bytes.NewBufferString(`{"filename":"big.mp4","content_type":"video/mp4","total_size":52428800}`)
	req, _ := http.NewRequest("POST", "/api/files/multipart/initiate", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["upload_id"] != "upload-abc" {
		t.Errorf("expected upload_id upload-abc, got %v", resp["upload_id"])
	}
}

func TestInitiateMultipartUpload_InvalidBody(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := bytes.NewBufferString(`{"content_type":"video/mp4"}`)
	req, _ := http.NewRequest("POST", "/api/files/multipart/initiate", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInitiateMultipartUpload_NoUser(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/files/multipart/initiate", handler.InitiateMultipartUpload)

	body := bytes.NewBufferString(`{"filename":"big.mp4","total_size":52428800}`)
	req, _ := http.NewRequest("POST", "/api/files/multipart/initiate", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUploadPart_Success(t *testing.T) {
	mockClient := &mockFileClient{
		uploadPartFunc: func(ctx context.Context, req *filev1.UploadPartRequest) (*filev1.UploadPartResponse, error) {
			if req.UploadId != "upload-abc" {
				t.Errorf("expected upload_id upload-abc, got %s", req.UploadId)
			}
			if req.PartNumber != 1 {
				t.Errorf("expected part_number 1, got %d", req.PartNumber)
			}
			return &filev1.UploadPartResponse{
				Etag:       "\"etag1\"",
				PartNumber: 1,
			}, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("chunk", "chunk.bin")
	part.Write([]byte("chunk data here"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files/multipart/upload-abc/part/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["etag"] != "\"etag1\"" {
		t.Errorf("expected etag \"etag1\", got %v", resp["etag"])
	}
}

func TestUploadPart_InvalidPartNumber(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("chunk", "chunk.bin")
	part.Write([]byte("data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files/multipart/upload-abc/part/notanumber", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUploadPart_MissingChunk(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("POST", "/api/files/multipart/upload-abc/part/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCompleteMultipartUpload_Success(t *testing.T) {
	mockClient := &mockFileClient{
		completeMultipartFunc: func(ctx context.Context, req *filev1.CompleteMultipartUploadRequest) (*filev1.FileResponse, error) {
			if req.UploadId != "upload-abc" {
				t.Errorf("expected upload_id upload-abc, got %s", req.UploadId)
			}
			if len(req.Parts) != 2 {
				t.Errorf("expected 2 parts, got %d", len(req.Parts))
			}
			return &filev1.FileResponse{
				File: &filev1.File{
					Id:          "file-done",
					Filename:    "big.mp4",
					Size:        52428800,
					ContentType: "video/mp4",
					CreatedAt:   1704067200,
				},
			}, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := bytes.NewBufferString(`{"parts":[{"part_number":1,"etag":"\"e1\""},{"part_number":2,"etag":"\"e2\""}]}`)
	req, _ := http.NewRequest("POST", "/api/files/multipart/upload-abc/complete", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	file := resp["file"].(map[string]interface{})
	if file["id"] != "file-done" {
		t.Errorf("expected file id file-done, got %v", file["id"])
	}
}

func TestCompleteMultipartUpload_InvalidBody(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := bytes.NewBufferString(`{}`)
	req, _ := http.NewRequest("POST", "/api/files/multipart/upload-abc/complete", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCompleteMultipartUpload_GRPCError(t *testing.T) {
	mockClient := &mockFileClient{
		completeMultipartFunc: func(ctx context.Context, req *filev1.CompleteMultipartUploadRequest) (*filev1.FileResponse, error) {
			return nil, status.Error(codes.NotFound, "upload session not found")
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := bytes.NewBufferString(`{"parts":[{"part_number":1,"etag":"\"e1\""}]}`)
	req, _ := http.NewRequest("POST", "/api/files/multipart/upload-abc/complete", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestAbortMultipartUpload_Success(t *testing.T) {
	mockClient := &mockFileClient{
		abortMultipartFunc: func(ctx context.Context, req *filev1.AbortMultipartUploadRequest) (*filev1.AbortMultipartUploadResponse, error) {
			if req.UploadId != "upload-abc" {
				t.Errorf("expected upload_id upload-abc, got %s", req.UploadId)
			}
			return &filev1.AbortMultipartUploadResponse{Success: true}, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("DELETE", "/api/files/multipart/upload-abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Errorf("expected success true, got %v", resp["success"])
	}
}

func TestAbortMultipartUpload_NoUser(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/api/files/multipart/:upload_id", handler.AbortMultipartUpload)

	req, _ := http.NewRequest("DELETE", "/api/files/multipart/upload-abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDownloadFile_Success(t *testing.T) {
	callCount := 0
	mockStream := &mockDownloadStream{
		recvFunc: func() (*filev1.DownloadFileResponse, error) {
			callCount++
			if callCount == 1 {
				return &filev1.DownloadFileResponse{
					Data: &filev1.DownloadFileResponse_Metadata{
						Metadata: &filev1.DownloadFileMetadata{
							Filename:    "test.txt",
							ContentType: "text/plain",
							Size:        12,
						},
					},
				}, nil
			} else if callCount == 2 {
				return &filev1.DownloadFileResponse{
					Data: &filev1.DownloadFileResponse_Chunk{
						Chunk: []byte("test content"),
					},
				}, nil
			}
			return nil, io.EOF
		},
	}

	mockClient := &mockFileClient{
		downloadFileFunc: func(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
			return mockStream, nil
		},
	}

	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/file-1/download", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("expected Content-Type text/plain, got %s", w.Header().Get("Content-Type"))
	}

	if w.Body.String() != "test content" {
		t.Errorf("expected body 'test content', got %s", w.Body.String())
	}
}

func TestMapGRPCError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"InvalidArgument", status.Error(codes.InvalidArgument, "bad"), http.StatusBadRequest},
		{"NotFound", status.Error(codes.NotFound, "missing"), http.StatusNotFound},
		{"PermissionDenied", status.Error(codes.PermissionDenied, "denied"), http.StatusForbidden},
		{"Unauthenticated", status.Error(codes.Unauthenticated, "unauth"), http.StatusUnauthorized},
		{"AlreadyExists", status.Error(codes.AlreadyExists, "dup"), http.StatusConflict},
		{"ResourceExhausted", status.Error(codes.ResourceExhausted, "limit"), http.StatusTooManyRequests},
		{"Unavailable", status.Error(codes.Unavailable, "down"), http.StatusServiceUnavailable},
		{"Internal", status.Error(codes.Internal, "err"), http.StatusInternalServerError},
		{"Unimplemented", status.Error(codes.Unimplemented, "not impl"), http.StatusInternalServerError},
		{"NonGRPC", errors.New("non-grpc error"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapGRPCError(tt.err)
			if got != tt.expected {
				t.Errorf("mapGRPCError(%v) = %d, want %d", tt.err, got, tt.expected)
			}
		})
	}
}

func TestDeleteFile_NoUser(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/api/files/:id", handler.DeleteFile)

	req, _ := http.NewRequest("DELETE", "/api/files/file-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDeleteFile_GRPCError(t *testing.T) {
	mockClient := &mockFileClient{
		deleteFileFunc: func(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
			return nil, status.Error(codes.NotFound, "file not found")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("DELETE", "/api/files/file-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteFile_PermissionDenied(t *testing.T) {
	mockClient := &mockFileClient{
		deleteFileFunc: func(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
			return nil, status.Error(codes.PermissionDenied, "not yours")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("DELETE", "/api/files/file-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestDeleteFile_InternalError(t *testing.T) {
	mockClient := &mockFileClient{
		deleteFileFunc: func(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
			return nil, status.Error(codes.Internal, "db error")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("DELETE", "/api/files/file-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUploadFile_NoUser(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/files", handler.UploadFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDownloadFile_NoUser(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/files/:id/download", handler.DownloadFile)

	req, _ := http.NewRequest("GET", "/api/files/f1/download", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDownloadFile_StreamError(t *testing.T) {
	mockClient := &mockFileClient{
		downloadFileFunc: func(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
			return nil, status.Error(codes.NotFound, "not found")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/f1/download", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetFile_NoUser(t *testing.T) {
	mockClient := &mockFileClient{}
	handler := NewFileHandler(mockClient)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/files/:id", handler.GetFile)

	req, _ := http.NewRequest("GET", "/api/files/file-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetFile_PermissionDenied(t *testing.T) {
	mockClient := &mockFileClient{
		getFileFunc: func(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
			return nil, status.Error(codes.PermissionDenied, "not authorized")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/file-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetFile_InternalError(t *testing.T) {
	mockClient := &mockFileClient{
		getFileFunc: func(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
			return nil, status.Error(codes.Internal, "db error")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/file-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetFile_InvalidArgument(t *testing.T) {
	mockClient := &mockFileClient{
		getFileFunc: func(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
			return nil, status.Error(codes.InvalidArgument, "invalid file id")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/bad-id", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUploadFile_SendError(t *testing.T) {
	mockStream := &mockUploadStream{
		sendFunc: func(req *filev1.UploadFileRequest) error {
			return status.Error(codes.Unavailable, "service down")
		},
	}
	mockClient := &mockFileClient{
		uploadFileFunc: func(ctx context.Context) (filev1.FileService_UploadFileClient, error) {
			return mockStream, nil
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUploadFile_AlreadyExists(t *testing.T) {
	mockStream := &mockUploadStream{
		sendFunc: func(req *filev1.UploadFileRequest) error { return nil },
		closeAndRecvFunc: func() (*filev1.FileResponse, error) {
			return nil, status.Error(codes.AlreadyExists, "file exists")
		},
	}
	mockClient := &mockFileClient{
		uploadFileFunc: func(ctx context.Context) (filev1.FileService_UploadFileClient, error) {
			return mockStream, nil
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestDownloadFile_RecvMetadataError(t *testing.T) {
	mockStream := &mockDownloadStream{
		recvFunc: func() (*filev1.DownloadFileResponse, error) {
			return nil, status.Error(codes.Internal, "stream error")
		},
	}
	mockClient := &mockFileClient{
		downloadFileFunc: func(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
			return mockStream, nil
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/f1/download", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestDownloadFile_PermissionDenied(t *testing.T) {
	mockClient := &mockFileClient{
		downloadFileFunc: func(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
			return nil, status.Error(codes.PermissionDenied, "not yours")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/f1/download", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestDownloadFile_InvalidArgument(t *testing.T) {
	mockClient := &mockFileClient{
		downloadFileFunc: func(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
			return nil, status.Error(codes.InvalidArgument, "invalid file id")
		},
	}
	handler := NewFileHandler(mockClient)
	router := setupFileTestRouter(handler)

	req, _ := http.NewRequest("GET", "/api/files/bad-id/download", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
