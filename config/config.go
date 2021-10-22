package config

import (
	"github.com/ducketlab/auth/client"
)

func newConfig() *Config {
	return &Config{
		App:  newDefaultApp(),
		Http: newDefaultHttp(),
		Grpc: newDefaultGrpc(),
		Log:  newDefaultLog(),
		Auth: newDefaultAuth(),
	}
}

type Config struct {
	App  *app  `toml:"app"`
	Http *http `toml:"http"`
	Grpc *grpc `toml:"grpc"`
	Log  *log  `toml:"log"`
	Auth *auth `toml:"auth"`
}

type app struct {
	Name string `toml:"name" env:"BOOK_APP_NAME"`
	Key  string `toml:"key" env:"BOOK_APP_NAME"`
}

func newDefaultApp() *app {
	return &app{
		Name: "book",
		Key:  "default",
	}
}

type http struct {
	Host string `toml:"host" env:"BOOK_HTTP_HOST"`
	Port string `toml:"port" env:"BOOK_HTTP_PORT"`
}

func (a *http) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultHttp() *http {
	return &http{
		Host: "127.0.0.1",
		Port: "8050",
	}
}

type grpc struct {
	Host string `toml:"host" env:"BOOK_GRPC_HOST"`
	Port string `toml:"port" env:"BOOK_GRPC_PORT"`
}

func (a *grpc) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultGrpc() *grpc {
	return &grpc{
		Host: "127.0.0.1",
		Port: "18050",
	}
}

type log struct {
	Level   string    `toml:"level" env:"BOOK_LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"BOOK_LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"BOOK_LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"BOOK_LOG_TO"`
}

func newDefaultLog() *log {
	return &log{
		Level:   "debug",
		PathDir: "logs",
		Format:  "text",
		To:      "stdout",
	}
}

type auth struct {
	Host         string `toml:"host" env:"BOOK_AUTH_HOST"`
	Port         string `toml:"port" env:"BOOK_AUTH_PORT"`
	ClientId     string `toml:"client_id" env:"BOOK_AUTH_CLIENT_ID"`
	ClientSecret string `toml:"client_secret" env:"BOOK_AUTH_CLIENT_SECRET"`
}

func (a *auth) Addr() string {
	return a.Host + ":" + a.Port
}

func (a *auth) Client() (*client.Client, error) {
	if client.C() == nil {
		config := client.NewDefaultConfig()

		config.SetAddress(a.Addr())
		config.SetClientCredentials(a.ClientId, a.ClientSecret)

		cli, err := client.NewClient(config)

		if err != nil {
			return nil, err
		}

		client.SetGlobal(cli)
	}

	return client.C(), nil
}

func newDefaultAuth() *auth {
	return &auth{}
}
