package server

import (
	pb "coauth/api/v1"
	v1 "coauth/api/v1"
	"coauth/internal/di"
	"coauth/internal/middleware"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthsvc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Start(deps *di.Container) error {
	listenAddr := deps.Config.Server.Grpc.Addr
	if len(listenAddr) == 0 {
		listenAddr = ":10103"
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		deps.Log.Error("failed to listen", zap.Error(err))
		return err
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
      middleware.TimeoutInterceptor(),
			middleware.RequestID(),
			middleware.LoggingInterceptor(deps.Log),
      middleware.ClientDisconnectInterceptor(),
			middleware.RecoverInterceptor(deps.Log),
		),
	)
	pb.RegisterCoauthServer(s, deps.Svc)
	// 健康检查
	healthSvc := healthsvc.NewServer()
	healthpb.RegisterHealthServer(s, healthSvc)
	updateHealthStatus(healthSvc, v1.Coauth_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	deps.Log.Info("starting grpc server", zap.String("addr", listenAddr))
	if err := s.Serve(listener); err != nil {
		deps.Log.Error("failed to serve", zap.Error(err))
		return err
	}
	return nil
}

func updateHealthStatus(
	h *healthsvc.Server,
	service string,
	status healthpb.HealthCheckResponse_ServingStatus) {
	h.SetServingStatus(service, status)
}
