package kinkku

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func StartUp() {
	fmt.Println(FgMagenta + banner + Reset)
	fmt.Println(FgCyan + Italic + slogan + Reset)
}

func GetArgs() {
	if len(os.Args) != 3 && len(os.Args) != 1 {
		fmt.Println(FgRed + "Oops, skill issue: You got the ham, but where is the mustard?" + Reset)
		fmt.Println("Wrong number of arguments.")
		fmt.Println("kinkku usage example:")
		fmt.Println("$ kinkku ./directory 6969")
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		path = "."
		port = "8080"
		fmt.Println()
		fmt.Println(FgCyan + "Using current directory and port 8080 by default." + Reset)
	} else {
		path = os.Args[1]
		port = os.Args[2]
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(FgRed + "Oops, skill issue: Can't find your fridge." + Reset)
		fmt.Println("Provided directory path does not exist.")
		fmt.Println("kinkku usage example:")
		fmt.Println("$ kinkku ./directory 6969")
		os.Exit(0)
	}
}

// Function to restart the Go server
func RestartServer() {
	// Kill the existing server process if it's running
	if serverCmd != nil && serverCmd.Process != nil {
		if err := serverCmd.Process.Kill(); err != nil {
			fmt.Println("Error killing server:", err)
		}
	}

	// Start the server again in a separate goroutine
	go func() {
		serverCmd = exec.Command("go", "run", ".")
		serverCmd.Dir = path // Set the working directory for the command
		serverCmd.Stdout = os.Stdout
		serverCmd.Stderr = os.Stderr
		if err := serverCmd.Start(); err != nil {
			fmt.Println(FgRed + "Oops, skill issue: there is no ham to serve." + Reset)
			fmt.Println("Error starting server:", err)
			os.Exit(0)
		}

		if restartCount == 0 {
			fmt.Println(FgGreen + "Server is up! Let's Go!" + Reset)
			restartCount++
		} else if restartCount == 69 {
			fmt.Println(FgGreen + "(" + strconv.Itoa(restartCount) + ") " + "Server restarted." + Reset)
			fmt.Println(FgCyan + string(noice) + Reset)
			restartCount++
		} else {
			fmt.Println(FgGreen + "(" + strconv.Itoa(restartCount) + ") " + "Server restarted." + Reset)
			restartCount++
		}
	}()
}

// Function to watch for file changes recursively in a directory
func WatchFiles(changes chan<- string) {
	fileModTimes := getFileModTimes()
	goFilesFound := false
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			goFilesFound = true // Set to true if Go file is found
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	if !goFilesFound {
		fmt.Println(FgRed + "Oops, skill issue: there is no ham to serve." + Reset)
		fmt.Println("No Go files found in provided directory.")
		os.Exit(0)
	}

	// Continuously monitor for file changes
	for {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			checkFileModifications(path, info, fileModTimes, changes)
			return nil
		})
		if err != nil {
			fmt.Println("Error:", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func checkFileModifications(path string, info os.FileInfo, fileModTimes map[string]time.Time, changes chan<- string) {
	modTime := info.ModTime()
	lastModTime := fileModTimes[path]
	if strings.HasSuffix(path, ".go") {
		if !lastModTime.IsZero() {
			if modTime.After(lastModTime) {
				fmt.Println(FgMagenta + "Go file modification detected:" + Reset + path)
				changes <- path
				fileModTimes[path] = modTime
				ModificationDetected = true // Set the flag indicating modification detected
			}
		} else {
			fmt.Println(FgYellow + "Go file found:" + Reset + path)
			changes <- path
			fileModTimes[path] = modTime
		}
	}
}

func getFileModTimes() map[string]time.Time {
	fileModTimes := make(map[string]time.Time)
	return fileModTimes
}
