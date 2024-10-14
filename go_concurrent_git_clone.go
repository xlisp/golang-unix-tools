package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func main() {
	repos := []string{
		"https://github.com/golang/go.git",
		"https://github.com/golang/tools.git",
		"https://github.com/golang/crypto.git",
		"https://github.com/golang/net.git",
		"https://github.com/golang/text.git",
		"https://github.com/golang/sync.git",
		"https://github.com/golang/sys.git",
		"https://github.com/golang/oauth2.git",
		"https://github.com/golang/mobile.git",
		"https://github.com/golang/image.git",
	}

	projectDir := "golang_projects"
	err := os.Mkdir(projectDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating project directory: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			cloneRepo(repo, projectDir)
		}(repo)
	}

	wg.Wait()
	fmt.Println("All repositories have been cloned.")
}

func cloneRepo(repo, projectDir string) {
	cmd := exec.Command("git", "clone", repo)
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error cloning %s: %v\n%s\n", repo, err, string(output))
	} else {
		fmt.Printf("Successfully cloned %s\n", repo)
	}
}
