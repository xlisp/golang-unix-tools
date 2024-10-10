package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Function to check if a file contains the given pattern
func containsPattern(filePath, pattern string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), pattern) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

// Function to walk through the directory and delete matching files
func main() {
	rootDir := "." // Specify the root directory you want to search in
	pattern := "# lib/jimw-code/" // The pattern to search for

	// Walk through all the files in the directory
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file contains the pattern
		matches, err := containsPattern(path, pattern)
		if err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			return nil
		}

		// If the file matches, delete it
		if matches {
			fmt.Printf("Deleting file: %s\n", path)
			err := os.Remove(path)
			if err != nil {
				log.Printf("Failed to delete file %s: %v", path, err)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
}
