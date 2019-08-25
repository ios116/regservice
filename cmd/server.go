package main

import (
	"fmt"
	"github.com/ios116/regservice/config"
	"github.com/ios116/regservice/server"
	"github.com/ios116/regservice/session"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfg := config.NewConfig()
	logger, _ := cfg.CreateLogger()
	defer logger.Sync()
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("cannot listen port", zap.Int("port", cfg.Port))
	}

	logger.Info("Server is starting", zap.String("host", cfg.Host), zap.Int("port", cfg.Port))

	grpcServer := grpc.NewServer()
	session.RegisterAuthCheckerServer(grpcServer, server.NewSessionManager(logger))
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Fatal("cannot start grpc server")
	}
}
