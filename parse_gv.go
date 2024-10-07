package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Edge struct {
	From string
	To   string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_graphviz_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	edges := parseGraphvizFile(filename)

	fmt.Println("Parsed edges:")
	for _, edge := range edges {
		fmt.Printf("%s -> %s\n", edge.From, edge.To)
	}
}

func parseGraphvizFile(filename string) []Edge {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var edges []Edge
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "->") {
			parts := strings.Split(line, "->")
			if len(parts) == 2 {
				from := strings.TrimSpace(parts[0])
				to := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))
				edges = append(edges, Edge{From: from, To: to})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return edges
}

