package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubRepo = "zanderson3j/o_and_m_online"
	updateCheckInterval = 24 * time.Hour
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type UpdateChecker struct {
	currentVersion string
	lastCheck      time.Time
}

func NewUpdateChecker() *UpdateChecker {
	return &UpdateChecker{
		currentVersion: "v" + currentVersion, // currentVersion from updater.go
	}
}

func (uc *UpdateChecker) CheckForUpdates() (bool, *GitHubRelease, error) {
	// Only check once per day
	if time.Since(uc.lastCheck) < updateCheckInterval {
		return false, nil, nil
	}
	uc.lastCheck = time.Now()

	// Get latest release from GitHub
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)
	resp, err := http.Get(url)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, nil, err
	}

	// Compare versions
	if release.TagName > uc.currentVersion {
		return true, &release, nil
	}

	return false, nil, nil
}

func (uc *UpdateChecker) DownloadAndInstall(release *GitHubRelease) error {
	// Find the right DMG for this architecture
	archName := "intel"
	if runtime.GOARCH == "arm64" {
		archName = "apple_silicon"
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, archName) && strings.HasSuffix(asset.Name, ".dmg") {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no DMG found for architecture %s", archName)
	}

	// Download to temp file
	tmpFile, err := os.CreateTemp("", "oam_update_*.dmg")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	log.Printf("Downloading update from %s", downloadURL)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return err
	}
	tmpFile.Close()

	// Mount the DMG
	log.Println("Mounting DMG...")
	mountCmd := exec.Command("hdiutil", "attach", tmpFile.Name(), "-nobrowse", "-quiet")
	output, err := mountCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to mount DMG: %w", err)
	}

	// Extract mount point
	lines := strings.Split(string(output), "\n")
	var mountPoint string
	for _, line := range lines {
		if strings.Contains(line, "/Volumes/") {
			parts := strings.Fields(line)
			mountPoint = parts[len(parts)-1]
			break
		}
	}

	if mountPoint == "" {
		return fmt.Errorf("could not find mount point")
	}

	defer func() {
		// Unmount DMG
		exec.Command("hdiutil", "detach", mountPoint, "-quiet").Run()
	}()

	// Copy new app to Applications
	appName := "O&M Game Room.app"
	srcPath := filepath.Join(mountPoint, appName)
	dstPath := filepath.Join("/Applications", appName)

	// Remove old app
	os.RemoveAll(dstPath)

	// Copy new app
	log.Printf("Installing update to %s", dstPath)
	copyCmd := exec.Command("cp", "-R", srcPath, dstPath)
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy app: %w", err)
	}

	log.Println("Update installed successfully!")
	
	// Restart the app
	log.Println("Restarting application...")
	exec.Command("open", dstPath).Start()
	os.Exit(0)

	return nil
}

func checkAndPromptForUpdate(gr *GameRoom) {
	checker := NewUpdateChecker()

	hasUpdate, release, err := checker.CheckForUpdates()
	if err != nil {
		log.Printf("Failed to check for updates: %v", err)
		return
	}

	if !hasUpdate {
		log.Println("No updates available - running latest version")
		return
	}

	log.Printf("New version available: %s (current: %s)", release.TagName, checker.currentVersion)

	// Set update notification state
	gr.updateAvailable = true
	gr.updateVersion = release.TagName
	gr.updateURL = fmt.Sprintf("https://github.com/%s/releases/tag/%s", githubRepo, release.TagName)

	log.Printf("Update available! Visit: %s", gr.updateURL)
}