package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection parameters
	connStr := "user=postgrest password=123456 dbname=text2code sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query to get the blog data
	rows, err := db.Query("SELECT id, name, content FROM blogs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the rows
	for rows.Next() {
		var id int
		var name, content string

		// Scan the row data into variables
		if err := rows.Scan(&id, &name, &content); err != nil {
			log.Fatal(err)
		}

		// Create markdown file with the id as the filename
		fileName := fmt.Sprintf("text2code_%d.md", id)
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Write markdown content to the file
		markdownContent := fmt.Sprintf("# %s\n\n%s", name, content)
		_, err = file.WriteString(markdownContent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Markdown file created: %s\n", fileName)
	}

	// Check for errors after iterating
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Markdown files exported successfully.")
}

