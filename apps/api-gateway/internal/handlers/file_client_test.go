package handlers

import (
	"context"
	"errors"
	"testing"

	filev1 "github.com/provsalt/DOP_P01_Team1/common/file/v1"
)

type mockFileClient struct {
	listFilesFunc         func(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error)
	getFileFunc           func(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error)
	deleteFileFunc        func(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error)
	uploadFileFunc        func(ctx context.Context) (filev1.FileService_UploadFileClient, error)
	downloadFileFunc      func(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error)
	initiateMultipartFunc func(ctx context.Context, req *filev1.InitiateMultipartUploadRequest) (*filev1.InitiateMultipartUploadResponse, error)
	uploadPartFunc        func(ctx context.Context, req *filev1.UploadPartRequest) (*filev1.UploadPartResponse, error)
	completeMultipartFunc func(ctx context.Context, req *filev1.CompleteMultipartUploadRequest) (*filev1.FileResponse, error)
	abortMultipartFunc    func(ctx context.Context, req *filev1.AbortMultipartUploadRequest) (*filev1.AbortMultipartUploadResponse, error)
}

func (m *mockFileClient) ListFiles(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error) {
	if m.listFilesFunc != nil {
		return m.listFilesFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) GetFile(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
	if m.getFileFunc != nil {
		return m.getFileFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) DeleteFile(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
	if m.deleteFileFunc != nil {
		return m.deleteFileFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) UploadFile(ctx context.Context) (filev1.FileService_UploadFileClient, error) {
	if m.uploadFileFunc != nil {
		return m.uploadFileFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) DownloadFile(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
	if m.downloadFileFunc != nil {
		return m.downloadFileFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) InitiateMultipartUpload(ctx context.Context, req *filev1.InitiateMultipartUploadRequest) (*filev1.InitiateMultipartUploadResponse, error) {
	if m.initiateMultipartFunc != nil {
		return m.initiateMultipartFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) UploadPart(ctx context.Context, req *filev1.UploadPartRequest) (*filev1.UploadPartResponse, error) {
	if m.uploadPartFunc != nil {
		return m.uploadPartFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) CompleteMultipartUpload(ctx context.Context, req *filev1.CompleteMultipartUploadRequest) (*filev1.FileResponse, error) {
	if m.completeMultipartFunc != nil {
		return m.completeMultipartFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) AbortMultipartUpload(ctx context.Context, req *filev1.AbortMultipartUploadRequest) (*filev1.AbortMultipartUploadResponse, error) {
	if m.abortMultipartFunc != nil {
		return m.abortMultipartFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileClient) Close() error {
	return nil
}

func TestFileServiceClientInterfaceImplemented(t *testing.T) {
	var _ FileServiceClient = (*mockFileClient)(nil)
}

func TestNewGRPCFileClient_ReturnsFileServiceClient(t *testing.T) {
	client, err := NewGRPCFileClient("localhost:8082")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if client == nil {
		t.Errorf("expected client, got nil")
	}

	if err := client.Close(); err != nil {
		t.Errorf("failed to close client: %v", err)
	}
}
