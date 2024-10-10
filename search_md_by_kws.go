package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Function to check if a file has a .md extension (Markdown file)
func isMarkdownFile(filePath string) bool {
	return filepath.Ext(filePath) == ".md"
}

// Function to search for keywords in a file
func searchKeywordsInFile(filePath string, keywords []string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if all keywords are found in the line
		foundAll := true
		for _, keyword := range keywords {
			if !strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				foundAll = false
				break
			}
		}
		if foundAll {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return false
	}

	return false
}

// Function to search through markdown files in the current directory
func searchMarkdownFiles(dir string, keywords []string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If it's a markdown file, search for the keywords
		if !info.IsDir() && isMarkdownFile(path) {
			if searchKeywordsInFile(path, keywords) {
				fmt.Printf("Keywords found in: %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %s\n", err)
	}
}

func main() {
	// Ensure keywords are provided as command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <keyword1> <keyword2> ...")
		return
	}

	// Command-line arguments after the program name are considered as keywords
	keywords := os.Args[1:]

	// Set the directory to the current path
	dir := "."

	// Search markdown files in the current directory for the keywords
	searchMarkdownFiles(dir, keywords)
}

