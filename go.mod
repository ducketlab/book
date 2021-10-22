module github.com/ducketlab/book

go 1.17

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/caarlos0/env/v6 v6.5.0
	github.com/ducketlab/auth v0.0.0-20211015124231-f2a1ff3e8137
	github.com/ducketlab/mingo v0.0.0-20211015124202-a23e3ffa85c4
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/spf13/cobra v1.2.1
	google.golang.org/grpc v1.38.0
)

require (
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/xid v1.3.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/ducketlab/mingo => ../mingo

replace github.com/ducketlab/auth => ../auth
