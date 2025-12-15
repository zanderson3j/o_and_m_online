package main

import (
	"log"
	"runtime"
)

const currentVersion = "1.0.17"

// checkForUpdates checks for new versions and logs the current version
func checkForUpdates(gr *GameRoom) {
	log.Printf("O&M Game Room v%s (%s/%s)", currentVersion, runtime.GOOS, runtime.GOARCH)

	// Check for updates in background
	go checkAndPromptForUpdate(gr)
}