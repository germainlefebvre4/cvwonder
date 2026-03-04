package cvserve

import "net"

type ServeInterface interface {
	StartLiveReloader(port int, outputDirectory string, inputFilePath string)
	StartServer(port int, outputDirectory string)
	StartServerOnListener(listener net.Listener, outputDirectory string)
	OpenBrowser(outputDirectory string, inputFilePath string)
}

type ServeServices struct{}

func NewServeServices() (ServeInterface, error) {
	return &ServeServices{}, nil
}
