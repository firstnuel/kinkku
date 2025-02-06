package main

import (
	"time"

	kinkku "github.com/firstnuel/kinkku/kinkku"
)

func main() {
	// Parse command-line arguments
	kinkku.GetArgs()

	// Display startup banner and slogan
	kinkku.StartUp()

	// Run the server for the first time
	kinkku.RestartServer()

	// Create a channel to receive file change events
	fileChanges := make(chan string)

	// Start watching for file changes in a separate goroutine
	go kinkku.WatchFiles(fileChanges)

	// Watch for file change events and restart the server
	for {
		select {
		case <-fileChanges:
			// If a modification has been detected, restart the server
			if kinkku.ModificationDetected {
				kinkku.RestartServer()
				kinkku.ModificationDetected = false // Reset the flag after restarting
			}
		default:
			// Avoid busy-waiting by adding a small sleep
			time.Sleep(100 * time.Millisecond)
		}
	}
}
