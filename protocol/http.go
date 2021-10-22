package protocol

import (
	"context"
	"fmt"
	"github.com/ducketlab/auth/client"
	"github.com/ducketlab/auth/pkg"
	"github.com/ducketlab/book/config"
	"github.com/ducketlab/mingo/http/middleware/recovery"
	"github.com/ducketlab/mingo/http/router"
	"github.com/ducketlab/mingo/http/router/httprouter"
	"github.com/ducketlab/mingo/logger"
	"github.com/ducketlab/mingo/logger/zap"
	"net/http"
	"time"
)

func NewHTTPService() *HttpService {
	r := httprouter.New()
	r.Use(recovery.NewWithLogger(zap.L().Named("recovery")))
	r.EnableApiRoot()

	server := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      25 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Addr:              config.C().Http.Addr(),
		Handler:           r,
	}
	return &HttpService{
		r:      r,
		server: server,
		l:      zap.L().Named("http-service"),
		c:      config.C(),
	}
}

type HttpService struct {
	r      router.Router
	l      logger.Logger
	c      *config.Config
	server *http.Server
}

func (s *HttpService) Start() error {

	if err := s.initGRPCClient(); err != nil {
		return err
	}

	if err := pkg.InitV1HttpApi(s.c.App.Name, s.r); err != nil {
		return err
	}

	s.l.Infof("http listen address: %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

func (s *HttpService) Stop() error {

	s.l.Info("start graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Errorf("graceful shutdown timeout, force exit")
	}

	return nil
}

func (s *HttpService) initGRPCClient() error {
	cf := client.NewDefaultConfig()
	cf.SetAddress(s.c.Grpc.Addr())
	cf.SetClientCredentials(s.c.Auth.ClientId, s.c.Auth.ClientSecret)
	cli, err := client.NewClient(cf)
	if err != nil {
		return err
	}
	client.SetGlobal(cli)
	return err
}
