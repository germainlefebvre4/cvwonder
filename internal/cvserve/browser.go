package cvserve

import (
	"cvwonder/internal/utils"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/sirupsen/logrus"
)

func OpenBrowser() {
	logrus.Debug("Opening browser")
	url := fmt.Sprintf("http://localhost:%d", utils.CliArgs.Port)
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
	utils.CheckError(err)
}
