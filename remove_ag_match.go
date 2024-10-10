package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Function to search for files with ag and delete them
func main() {
	// Define the pattern to search for
	pattern := "# lib/jimw-code/"
	
	// Run the 'ag' command to search for files containing the pattern
	cmd := exec.Command("ag", "-l", pattern)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("ag command failed: %v", err)
	}

	// Get the list of files from the output
	files := strings.Split(out.String(), "\n")
	
	// Iterate through each file and delete it
	for _, file := range files {
		if file == "" {
			continue // Skip empty strings
		}
		// Delete the file
		err := os.Remove(file)
		if err != nil {
			log.Printf("Failed to delete file %s: %v", file, err)
		} else {
			fmt.Printf("Deleted file: %s\n", file)
		}
	}
}

