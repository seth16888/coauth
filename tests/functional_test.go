package tests

import (
	v1 "coauth/api/v1"
	"coauth/internal/di"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer(configFile string) (*grpc.Server, *bufconn.Listener) {
	appDeps := di.NewContainer(configFile)

	listener := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s, appDeps.Svc)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return s, listener
}

func registerServices(s *grpc.Server, svc v1.CoauthServer) {
	// 注册服务
	v1.RegisterCoauthServer(s, svc)
}

func TestAuthorize2(t *testing.T) {
	configFile := ""
	if os.Getenv("COAUTH_CONF") != "" {
		configFile = os.Getenv("COAUTH_CONF")
	}
	if configFile == "" {
		t.Fatalf("COAUTH_CONF is not set")
	}

	s, listener := startTestGrpcServer(configFile)
	defer s.GracefulStop()

	time.Sleep(2 * time.Second)

	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return listener.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use grpc.DialContext to establish a connection
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufconnDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	coauthClient := v1.NewCoauthClient(conn)
	req := &v1.AuthorizeRequest{
		ClientId:     "1",
		RedirectUri:  "http://www.seth.vip:10200/oauth/callback",
		ResponseType: "authentication_code",
		State:        "kxjkjds",
	}
	resp, err := coauthClient.Authorize(context.Background(), req)
	if err != nil {
		t.Fatalf("could not authorize: %v", err)
	}
	if resp.Code == "" {
		t.Fatalf("expected code, got %v", resp.Code)
	}
	if resp.State != req.State {
		t.Fatalf("expected state %v, got %v", req.State, resp.State)
	}
	if !strings.HasPrefix(resp.RedirectUri, req.RedirectUri) {
		t.Fatalf("expected redirect_uri %v, got %v", req.RedirectUri, resp.RedirectUri)
	}
}
