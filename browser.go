package main

import (
	"log"
	"os/exec"
	"runtime"
)

// OpenBrowser opens the specified URL in the user's default browser
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Printf("Unsupported platform for opening browser: %s", runtime.GOOS)
		return
	}

	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	} else {
		log.Printf("Opened browser to: %s", url)
	}
}
