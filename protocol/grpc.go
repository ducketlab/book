package protocol

import (
	"context"
	"fmt"
	"github.com/ducketlab/auth/pkg/endpoint"
	"github.com/ducketlab/book/config"
	"github.com/ducketlab/book/pkg"
	"github.com/ducketlab/book/version"
	"github.com/ducketlab/mingo/logger"
	"github.com/ducketlab/mingo/logger/zap"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"net"
)

func NewGRPCService(interceptors ...grpc.UnaryServerInterceptor) *GrpcService {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		interceptors...,
	)))

	return &GrpcService{
		svr: grpcServer,
		l:   zap.L().Named("grpc-service"),
		c:   config.C(),
	}
}

type GrpcService struct {
	svr *grpc.Server
	l   logger.Logger
	c   *config.Config
}

func (s *GrpcService) Start() error {

	// Load all grpc service
	pkg.InitV1GrpcApi(s.svr)

	lis, err := net.Listen("tcp", s.c.Grpc.Addr())
	if err != nil {
		return err
	}

	s.l.Infof("grpc listen address: %s", s.c.Grpc.Addr())

	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start grpc service error, %s", err.Error())
	}
	return nil
}

func (s *GrpcService) Stop() error {
	s.l.Info("start graceful shutdown")

	s.svr.GracefulStop()

	return nil
}

func (s *GrpcService) RegistryEndpoints() error {
	cli, err := s.c.Auth.Client()

	if err != nil {
		return nil
	}

	req := endpoint.NewRegistryRequest(version.Short(), pkg.HttpEntry().Items)

	_, err = cli.Endpoint().Registry(context.Background(), req)

	return err
}
