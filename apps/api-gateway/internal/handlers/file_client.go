package handlers

import (
	"context"

	filev1 "github.com/provsalt/DOP_P01_Team1/common/file/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FileServiceClient interface {
	ListFiles(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error)
	GetFile(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error)
	DeleteFile(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error)
	UploadFile(ctx context.Context) (filev1.FileService_UploadFileClient, error)
	DownloadFile(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error)
	Close() error
}

type grpcFileClient struct {
	conn   *grpc.ClientConn
	client filev1.FileServiceClient
}

func NewGRPCFileClient(addr string) (FileServiceClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcFileClient{
		conn:   conn,
		client: filev1.NewFileServiceClient(conn),
	}, nil
}

func (c *grpcFileClient) ListFiles(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error) {
	return c.client.ListFiles(ctx, req)
}

func (c *grpcFileClient) GetFile(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
	return c.client.GetFile(ctx, req)
}

func (c *grpcFileClient) DeleteFile(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
	return c.client.DeleteFile(ctx, req)
}

func (c *grpcFileClient) UploadFile(ctx context.Context) (filev1.FileService_UploadFileClient, error) {
	return c.client.UploadFile(ctx)
}

func (c *grpcFileClient) DownloadFile(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
	return c.client.DownloadFile(ctx, req)
}

func (c *grpcFileClient) Close() error {
	return c.conn.Close()
}
