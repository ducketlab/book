package pkg

import (
	"github.com/ducketlab/mingo/pb/http"
	"google.golang.org/grpc"
)

var (
	servers       []Service
	successLoaded []string
	entrySet      = http.NewEntrySet()
)

type Service interface {
	Config() error
	HttpEntry() *http.EntrySet
}

func HttpEntry() *http.EntrySet {
	return entrySet
}

func InitV1GrpcApi(server *grpc.Server) {

}

func GetGrpcPathEntry(path string) *http.Entry {
	es := HttpEntry()

	for i := range es.Items {
		if es.Items[i].Path == path {
			return es.Items[i]
		}
	}

	return nil
}

func addService(name string, svr Service) {
	servers = append(servers, svr)
	successLoaded = append(successLoaded, name)
}

func LoadedService() []string {
	return successLoaded
}

func InitService() error {
	for _, s := range servers {
		if err := s.Config(); err != nil {
			return err
		}
		entrySet.Merge(s.HttpEntry())
	}
	return nil
}
