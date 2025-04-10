package server

import (
	pb "coauth/api/v1"
	"coauth/internal/di"
	"coauth/internal/middleware"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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
			middleware.RequestID(),
			middleware.LoggingInterceptor(deps.Log),
		),
	)
	pb.RegisterCoauthServer(s, deps.Svc)

	deps.Log.Info("starting grpc server", zap.String("addr", listenAddr))
	if err := s.Serve(listener); err != nil {
		deps.Log.Error("failed to serve", zap.Error(err))
		return err
	}
	return nil
}
