package simplecli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Init() {
	dirPath := flag.String("path", ".", "Directory path to list files")
	extension := flag.String("ext", "", "File extension to filter by (e.g., .txt)")

	flag.Parse()

	var files []string
	var err error

	if *extension != "" {
		files, err = FilterFilesByExtension(*dirPath, *extension)
	} else {
		files, err = ListFiles(*dirPath)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, file := range files {
		fmt.Println(file)
	}
}

func ListFiles(directory string) ([]string, error) {
	var fileList []string

	files, err := os.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileList = append(fileList, file.Name())
	}

	return fileList, nil
}

func FilterFilesByExtension(directory, extension string) ([]string, error) {
	var filteredFiles []string

	files, err := os.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	// Normalize the extension to ensure it starts with a dot
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == extension {
			filteredFiles = append(filteredFiles, file.Name())
		}
	}
	return filteredFiles, nil
}
