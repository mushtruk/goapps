package simplecli

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CLIOptions struct {
	DirPath   string
	Extension string
	MaxSize   int
	Recursive bool
	Output    string
}

func ParseFlags() CLIOptions {
	dirPath := flag.String("path", ".", "Directory path to list files")
	extension := flag.String("ext", "", "File extension to filter by (e.g., .txt)")
	size := flag.Int("size", -1, "File size to filter by")
	recursive := flag.Bool("recursive", false, "List files in directories recursively")
	output := flag.String("output", "text", "Output format (e.g., 'text', 'JSON'")

	flag.Parse()

	return CLIOptions{
		DirPath:   *dirPath,
		Extension: *extension,
		MaxSize:   *size,
		Recursive: *recursive,
		Output:    *output,
	}
}

func Init(opts CLIOptions) {
	var files []string
	var err error

	// List files, either recursively or not based on the opts.Recursive flag
	if opts.Recursive {
		files, err = ListFilesRecursively(opts.DirPath)
	} else {
		files, err = ListFiles(opts.DirPath)
	}

	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	// Apply the extension filter if specified
	if opts.Extension != "" {
		files, err = FilterFilesByExtension(files, opts.Extension)
		if err != nil {
			log.Fatalf("Error filtering files by extension: %v", err)
		}
	}

	if opts.MaxSize >= 0 {
		files, err = FilterFilesBySize(files, opts.MaxSize)
		if err != nil {
			log.Fatalf("Error filtering files by size: %v", err)
		}
	}

	// Handle different output formats
	switch opts.Output {
	case "json":
		jsonData, err := json.Marshal(files)
		if err != nil {
			log.Fatalf("Error marshaling files to JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	default:
		for _, file := range files {
			fmt.Println(file)
		}
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

func FilterFilesByExtension(files []string, extension string) ([]string, error) {
	var filteredFiles []string
	for _, file := range files {
		if strings.HasSuffix(file, extension) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

func FilterFilesBySize(filePaths []string, maxSize int) ([]string, error) {
	var filteredFiles []string

	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		if info.Size() <= int64(maxSize) {
			filteredFiles = append(filteredFiles, filePath)
		}
	}

	return filteredFiles, nil
}

func ListFilesRecursively(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(directory, path)
			if err != nil {
				return err
			}
			files = append(files, relativePath)
		}
		return nil
	})

	return files, err
}
