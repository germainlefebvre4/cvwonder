package cvserve

import "net"

type ServeInterface interface {
	StartLiveReloader(port int, outputDirectory string, inputFilePath string) (<-chan error, error)
	StartServer(port int, outputDirectory string) (<-chan error, error)
	StartServerOnListener(listener net.Listener, outputDirectory string, ready chan<- struct{})
	OpenBrowser(outputDirectory string, inputFilePath string) error
}

type ServeServices struct{}

func NewServeServices() (ServeInterface, error) {
	return &ServeServices{}, nil
}
