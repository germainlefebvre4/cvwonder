package watcher

import (
	"github.com/germainlefebvre4/cvwonder/internal/cvparser"
	"github.com/germainlefebvre4/cvwonder/internal/cvrender"
)

type WatcherInterface interface {
	ObserveFileEvents(baseDirectory string, outputDirectory string, inputFilePath string, themeName string, exportFormat string)
}

type WatcherServices struct {
	ParserService cvparser.ParserInterface
	RenderService cvrender.RenderInterface
}

func NewWatcherServices(
	parserService cvparser.ParserInterface,
	renderService cvrender.RenderInterface,
) (WatcherInterface, error) {
	return &WatcherServices{
		ParserService: parserService,
		RenderService: renderService,
	}, nil
}
