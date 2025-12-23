package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// HotReload watches files and restarts the application on changes
type HotReload struct {
	ProjectPath string
	Command     string
	Args        []string
	cmd         *exec.Cmd
	restart     chan bool
}

// NewHotReload creates a new hot reload watcher
func NewHotReload(projectPath string) *HotReload {
	return &HotReload{
		ProjectPath: projectPath,
		Command:     "go",
		Args:        []string{"run", "main.go"},
		restart:     make(chan bool, 1),
	}
}

// Start begins watching and running the application
func (hr *HotReload) Start() error {
	fmt.Println("🔥 Hot reload enabled - watching for changes...")
	fmt.Println("   Watching: *.go, *.yaml, *.yml, *.html")
	fmt.Println()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start initial process
	hr.startProcess()

	// Watch for file changes
	go hr.watchFiles()

	// Main loop
	for {
		select {
		case <-hr.restart:
			hr.stopProcess()
			time.Sleep(100 * time.Millisecond)
			hr.startProcess()
		case <-sigChan:
			fmt.Println("\n⏹  Shutting down...")
			hr.stopProcess()
			return nil
		}
	}
}

// startProcess starts the application process
func (hr *HotReload) startProcess() {
	fmt.Printf("🚀 Starting application...\n")

	hr.cmd = exec.Command(hr.Command, hr.Args...)
	hr.cmd.Dir = hr.ProjectPath
	hr.cmd.Stdout = os.Stdout
	hr.cmd.Stderr = os.Stderr

	if err := hr.cmd.Start(); err != nil {
		log.Printf("Failed to start: %v", err)
		return
	}

	fmt.Printf("✓ Running (PID: %d)\n\n", hr.cmd.Process.Pid)
}

// stopProcess stops the application process
func (hr *HotReload) stopProcess() {
	if hr.cmd != nil && hr.cmd.Process != nil {
		fmt.Println("⏸  Stopping application...")
		hr.cmd.Process.Kill()
		hr.cmd.Wait()
	}
}

// watchFiles monitors file changes
func (hr *HotReload) watchFiles() {
	lastMod := make(map[string]time.Time)
	
	// Get initial file states
	filepath.Walk(hr.ProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if hr.shouldWatch(path) {
			lastMod[path] = info.ModTime()
		}
		return nil
	})

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		changed := false
		
		filepath.Walk(hr.ProjectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !hr.shouldWatch(path) {
				return nil
			}

			// Check if file is new or modified
			if lastModTime, exists := lastMod[path]; !exists || info.ModTime().After(lastModTime) {
				lastMod[path] = info.ModTime()
				if exists { // Only trigger on modification, not initial scan
					fmt.Printf("📝 File changed: %s\n", filepath.Base(path))
					changed = true
				}
			}
			return nil
		})

		if changed {
			select {
			case hr.restart <- true:
			default:
			}
		}
	}
}

// shouldWatch determines if a file should trigger a reload
func (hr *HotReload) shouldWatch(path string) bool {
	// Skip directories
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return false
	}

	// Skip hidden files and directories
	if len(path) > 0 && path[0] == '.' {
		return false
	}

	// Skip tmp and vendor directories
	if filepath.Base(filepath.Dir(path)) == "tmp" || 
	   filepath.Base(filepath.Dir(path)) == "vendor" ||
	   filepath.Base(filepath.Dir(path)) == ".git" {
		return false
	}

	// Watch specific extensions
	ext := filepath.Ext(path)
	watchExtensions := map[string]bool{
		".go":   true,
		".yaml": true,
		".yml":  true,
		".html": true,
		".tmpl": true,
	}

	return watchExtensions[ext]
}
