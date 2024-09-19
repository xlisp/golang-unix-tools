package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

var functionCalls = make(map[string][]string)

// Parse the function calls in a function declaration
func inspectFuncDecl(node ast.Node) bool {
	// We are only interested in function declarations
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return true
	}

	funcName := funcDecl.Name.Name

	// Traverse the function body to find function calls
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Get the name of the called function
		switch fun := callExpr.Fun.(type) {
		case *ast.Ident:
			// Simple function call
			functionCalls[funcName] = append(functionCalls[funcName], fun.Name)
		case *ast.SelectorExpr:
			// Method call (e.g., obj.Method)
			functionCalls[funcName] = append(functionCalls[funcName], fun.Sel.Name)
		}
		return true
	})

	return true
}

// Generate the Graphviz DOT format output
func generateDot() {
	fmt.Println("digraph G {")
	for caller, callees := range functionCalls {
		for _, callee := range callees {
			fmt.Printf("    \"%s\" -> \"%s\";\n", caller, callee)
		}
	}
	fmt.Println("}")
}

// Parse all Go files in the given directory
func parseGoFilesInDir(dir string) {
	fs := token.NewFileSet()

	// Walk through the directory to find Go files
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .go files
		if filepath.Ext(path) == ".go" {
			node, err := parser.ParseFile(fs, path, nil, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
				return err
			}

			// Walk through the AST and inspect function declarations
			ast.Inspect(node, inspectFuncDecl)
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking through directory: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_directory>")
		return
	}

	// Parse the directory
	dir := os.Args[1]
	parseGoFilesInDir(dir)

	// Generate Graphviz DOT output
	generateDot()
}

// run : is perfect！=》go run show_fun_refs_project.go /Users/emacspy/GoPro/xxxxx

