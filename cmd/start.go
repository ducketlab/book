package cmd

import (
	"errors"
	"fmt"
	"github.com/ducketlab/auth/client"
	"github.com/ducketlab/book/config"
	"github.com/ducketlab/book/pkg"
	"github.com/ducketlab/book/protocol"
	"github.com/ducketlab/mingo/logger"
	"github.com/ducketlab/mingo/logger/zap"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	confType string
	confFile string
)

var serviceCmd = &cobra.Command{
	Use: "start",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := loadGlobalConfig(confType); err != nil {
			return err
		}

		if err := loadGlobalLogger(); err != nil {
			return err
		}

		if err := pkg.InitService(); err != nil {
			return err
		}

		conf := config.C()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

		svr, err := newService(conf)
		if err != nil {
			return err
		}

		go svr.waitSign(ch)

		if err := svr.start(); err != nil {
			if !strings.Contains(err.Error(), "http: Server closed") {
				return err
			}
		}

		return nil
	},
}

func loadGlobalConfig(configType string) error {
	switch configType {
	case "file":
		err := config.LoadConfigFromToml(confFile)
		if err != nil {
			return err
		}
	case "env":
		err := config.LoadConfigFromEnv()
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown config type")
	}

	return nil
}

func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)

	lc := config.C().Log
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}

	zapConfig := zap.DefaultConfig()
	zapConfig.Level = level

	switch lc.To {
	case config.ToStdout:
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case config.ToFile:
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = lc.PathDir
	}

	switch lc.Format {
	case config.JSONFormat:
		zapConfig.Json = true
	}

	if err := zap.Configure(zapConfig); err != nil {
		return err
	}

	zap.L().Named("init").Info(logInitMsg)
	return nil
}

func newService(cnf *config.Config) (*service, error) {
	cli, err := cnf.Auth.Client()

	if err != nil {
		return nil, err
	}

	auther := client.NewGrpcAuthAuther(pkg.GetGrpcPathEntry, cli)

	auther.SetLogger(zap.L().Named("grpc-auther"))

	grpc := protocol.NewGRPCService(auther.AuthUnaryServerInterceptor())
	http := protocol.NewHTTPService()

	svr := &service{
		grpc: grpc,
		http: http,
		log:  zap.L().Named("cli"),
	}

	return svr, nil
}

func (s *service) start() error {

	s.log.Infof("loaded domain pkg: %v", pkg.LoadedService())

	s.log.Infof("loaded http service: %s", pkg.LoadedHttp())

	s.log.Info("start registry endpoints ...")
	if err := s.grpc.RegistryEndpoints(); err != nil {
		s.log.Warnf("registry endpoints error, %s", err)
	} else {
		s.log.Debug("service endpoints registry success")
	}

	go s.grpc.Start()
	return s.http.Start()
}

func (s *service) waitSign(sign chan os.Signal) {
	for {
		select {
		case sg := <-sign:
			switch v := sg.(type) {
			default:
				s.log.Infof("receive signal '%v', start graceful shutdown", v.String())
				if err := s.grpc.Stop(); err != nil {
					s.log.Errorf("grpc graceful shutdown err: %s, force exit", err)
				}
				s.log.Info("grpc service stop complete")
				if err := s.http.Stop(); err != nil {
					s.log.Errorf("http graceful shutdown err: %s, force exit", err)
				}
				s.log.Infof("http service stop complete")
				return
			}
		}
	}
}

type service struct {
	http *protocol.HttpService
	grpc *protocol.GrpcService

	log logger.Logger
}

func init() {
	serviceCmd.Flags().StringVarP(&confType, "config-type",
		"t", "file", "the service config type [file/env/etcd]")

	serviceCmd.Flags().StringVarP(&confFile, "config-file",
		"f", "etc/book.toml", "the service config from file")

	RootCmd.AddCommand(serviceCmd)
}
