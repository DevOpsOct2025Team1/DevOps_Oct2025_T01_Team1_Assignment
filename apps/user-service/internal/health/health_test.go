package health

import (
	"context"
	"testing"

	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

type watchServerStub struct {
	sent []*grpc_health_v1.HealthCheckResponse
}

func (w *watchServerStub) Send(resp *grpc_health_v1.HealthCheckResponse) error {
	w.sent = append(w.sent, resp)
	return nil
}

func (w *watchServerStub) SetHeader(metadata.MD) error  { return nil }
func (w *watchServerStub) SendHeader(metadata.MD) error { return nil }
func (w *watchServerStub) SetTrailer(metadata.MD)       {}
func (w *watchServerStub) Context() context.Context     { return context.Background() }
func (w *watchServerStub) SendMsg(any) error            { return nil }
func (w *watchServerStub) RecvMsg(any) error            { return nil }

func TestHealthServer_Check_ReturnsServing(t *testing.T) {
	srv := NewHealthServer()
	resp, err := srv.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("expected SERVING, got %v", resp.GetStatus())
	}
}

func TestHealthServer_Watch_SendsServing(t *testing.T) {
	srv := NewHealthServer()
	ws := &watchServerStub{}

	if err := srv.Watch(&grpc_health_v1.HealthCheckRequest{}, ws); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ws.sent) != 1 {
		t.Fatalf("expected 1 message sent, got %d", len(ws.sent))
	}
	if ws.sent[0].GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("expected SERVING, got %v", ws.sent[0].GetStatus())
	}
}
