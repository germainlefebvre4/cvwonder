package cvserve

import (
	"fmt"
	"log/slog"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
)

func (s *ServeServices) OpenBrowser(outputDirectory string, inputFilePath string) error {
	slog.Debug("Opening browser")

	// Input file
	inputFilenameExt := path.Base(inputFilePath)
	inputFilename := inputFilenameExt[:len(inputFilenameExt)-len(path.Ext(inputFilenameExt))]

	// Output file
	outputFilename := filepath.Base(inputFilename) + ".html"

	url := fmt.Sprintf("http://localhost:%d/%s", utils.CliArgs.Port, outputFilename)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return fmt.Errorf("opening browser at %s: %w", url, err)
	}
	return nil
}
