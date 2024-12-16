package serve

import (
	"fmt"
	"os/exec"
	"rendercv/internal/utils"
	"runtime"
)

func OpenBrowser(url string) {
	fmt.Println("Opening browser")
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
