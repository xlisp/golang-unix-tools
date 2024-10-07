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
	edges := []Edge{
		{"a", "b"},
		{"b", "c"},
		{"a", "c"},
	}

	writeGraphvizFile("graph.gv", edges)
	parsedEdges := parseGraphvizFile("graph.gv")

	fmt.Println("Parsed edges:")
	for _, edge := range parsedEdges {
		fmt.Printf("%s -> %s\n", edge.From, edge.To)
	}
}

func writeGraphvizFile(filename string, edges []Edge) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = writer.WriteString("digraph {\n")
	if err != nil {
		return err
	}

	for _, edge := range edges {
		_, err = writer.WriteString(fmt.Sprintf("  %s -> %s;\n", edge.From, edge.To))
		if err != nil {
			return err
		}
	}

	_, err = writer.WriteString("}\n")
	return err
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
