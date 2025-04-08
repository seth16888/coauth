package tests

import (
	v1 "coauth/api/v1"
	"context"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 测试/v1/authorize
func TestAuthorize(t *testing.T) {
	req := v1.AuthorizeRequest{
		ClientId:     "1",
		RedirectUri:  "http://www.seth.vip:10200/oauth/callback",
		ResponseType: "authentication_code",
		State:        "kxjkjds",
	}

	// Bug 修复：使用 grpc.Dial 来建立连接，同时使用 grpc.WithTransportCredentials 来设置传输凭证
	conn, err := grpc.NewClient("localhost:10101",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	client := v1.NewCoauthClient(conn)
	ctx := context.Background()
	resp, err := client.Authorize(ctx, &req)
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
