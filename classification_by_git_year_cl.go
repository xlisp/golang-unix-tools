package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run script.go <directory>")
		return
	}

	dir := os.Args[1]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			year, err := getGitCommitYear(filepath.Join(dir, file.Name()))
			if err != nil {
				fmt.Printf("Error getting commit year for %s: %v\n", file.Name(), err)
				continue
			}

			yearDir := filepath.Join(dir, year)
			err = os.MkdirAll(yearDir, 0755)
			if err != nil {
				fmt.Printf("Error creating year directory %s: %v\n", yearDir, err)
				continue
			}

			oldPath := filepath.Join(dir, file.Name())
			newPath := filepath.Join(yearDir, file.Name())
			err = os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Printf("Error moving file %s: %v\n", file.Name(), err)
				continue
			}

			fmt.Printf("Moved %s to %s\n", file.Name(), newPath)
		}
	}
}

func getGitCommitYear(filePath string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=short", "--", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	dateStr := strings.TrimSpace(string(output))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", date.Year()), nil
}
