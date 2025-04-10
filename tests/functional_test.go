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

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthsvc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var (
	grpcServer *grpc.Server
	listener   *bufconn.Listener
	configFile string
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
	healthpb.RegisterHealthServer(s, healthsvc.NewServer())

}

func createGrpcConnection(t *testing.T) *grpc.ClientConn {
	// 定义一个 bufconn 拨号器，用于在测试中建立连接
	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return listener.Dial()
	}

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 使用 grpc.DialContext 建立连接
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufconnDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	return conn
}

func TestMain(m *testing.M) {
	configFile = ""
	envPath := os.Getenv("COAUTH_CONF")
	if envPath != "" {
		configFile = envPath
	}
	if configFile == "" {
		log.Fatalf("Env COAUTH_CONF is not set")
	}

	grpcServer, listener = startTestGrpcServer(configFile)
	time.Sleep(2 * time.Second)

	// 运行所有测试用例
	code := m.Run()

	// 清理资源
	grpcServer.GracefulStop()
	listener.Close()

	os.Exit(code)
}

func TestAuthorize2(t *testing.T) {
	conn := createGrpcConnection(t)
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

func TestCaptcha(t *testing.T) {
	conn := createGrpcConnection(t)
	defer conn.Close()

	Convey("TestCaptcha", t, func() {
		coauthClient := v1.NewCoauthClient(conn)
		req := &v1.CaptchaRequest{}
		resp, err := coauthClient.Captcha(context.Background(), req)

		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.CaptchaKey, ShouldNotBeEmpty)
		So(resp.CaptchaValue, ShouldNotBeEmpty)

	})
}

func TestPasswordLogin(t *testing.T) {
	conn := createGrpcConnection(t)
	defer conn.Close()

	Convey("When request params is invalid", t, func() {
		client := v1.NewCoauthClient(conn)
		req := &v1.LoginRequest{
			Username:    "",
			Password:    "",
			CaptchaKey:  "",
			CaptchaCode: "",
		}

		resp, err := client.Login(context.Background(), req)
		So(err, ShouldNotBeNil)
		So(resp, ShouldBeNil)

		st := status.Convert(err)
		So(st.Code(), ShouldEqual, codes.InvalidArgument)
	})

	Convey("When database no data", t, func() {
		client := v1.NewCoauthClient(conn)
		req := &v1.LoginRequest{
			Username:    "admin",
			Password:    "123456",
			CaptchaKey:  "123412341234",
			CaptchaCode: "1234",
		}

		resp, err := client.Login(context.Background(), req)
		So(err, ShouldNotBeNil)
		So(resp, ShouldBeNil)
		st := status.Convert(err)
		So(st.Code(), ShouldEqual, 10400)
		So(st.Message(), ShouldContainSubstring, "invalid username or password")
	})

	Convey("When request params is valid", t, func() {
		client := v1.NewCoauthClient(conn)
		req := &v1.LoginRequest{
			Username:    "admin1",
			Password:    "admin1",
			CaptchaKey:  "123412341234",
			CaptchaCode: "1234",
		}

		resp, err := client.Login(context.Background(), req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.AccessToken, ShouldNotBeEmpty)
		So(resp.RefreshToken, ShouldNotBeEmpty)
		So(resp.TokenType, ShouldEqual, "Bearer")
		So(resp.ExpiresIn, ShouldBeGreaterThan, 0)
	})
}

// HealthCheck 测试健康检查
func TestHealthCheck(t *testing.T) {
	conn := createGrpcConnection(t)
	defer conn.Close()

  Convey("When service alive", t, func() {
		client := healthpb.NewHealthClient(conn)
		req := healthpb.HealthCheckRequest{
		}
		resp, err := client.Check(context.Background(), &req)

		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.Status, ShouldEqual, healthpb.HealthCheckResponse_SERVING)
	})
}
