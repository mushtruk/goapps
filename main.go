package main

import (
	"flag"
	"fmt"
	"goapps/simplecli"
	"log"
)

func main() {
	dirPath := flag.String("path", ".", "Directory path to list files")
	flag.Parse()

	// Call the ListFiles function with the provided directory path
	files, err := simplecli.ListFiles(*dirPath)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	// Print the list of files
	for _, file := range files {
		fmt.Println(file)
	}
}
