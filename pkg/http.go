package pkg

import (
	"github.com/ducketlab/mingo/http/router"
	"strings"
)

var (
	httpV1ApiMap = make(map[string]HttpApi)
)

func LoadedHttp() []string {
	var apis []string
	for k := range httpV1ApiMap {
		apis = append(apis, k)
	}
	return apis
}

type HttpApi interface {
	Registry(router.SubRouter)
	Config() error
}

func RegistryV1Http(name string, api HttpApi) {
	if _, ok := httpV1ApiMap[name]; ok {
		panic("http api " + name + " has registry")
	}
	httpV1ApiMap[name] = api
}

func InitV1HttpApi(pathPrefix string, root router.Router) error {
	for _, api := range httpV1ApiMap {
		if err := api.Config(); err != nil {
			return err
		}

		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}

		api.Registry(root.SubRouter(pathPrefix + "/v1"))
	}

	return nil
}
