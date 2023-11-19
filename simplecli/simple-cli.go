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

type FilterFunc func([]string) ([]string, error)

func Init(opts CLIOptions) {
	files, err := listFilesBasedOnFlag(opts)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	filters := getFilters(opts)
	for _, filter := range filters {
		files, err = filter(files)
		if err != nil {
			log.Fatalf("Error applying filter: %v", err)
		}
	}

	err = outputFiles(files, opts.Output)
	if err != nil {
		log.Fatalf("Error outputting files: %v", err)
	}
}

func listFilesBasedOnFlag(opts CLIOptions) ([]string, error) {
	if opts.Recursive {
		return ListFilesRecursively(opts.DirPath)
	}
	return ListFiles(opts.DirPath)
}

func getFilters(opts CLIOptions) []FilterFunc {
	var filters []FilterFunc

	if opts.Extension != "" {
		filters = append(filters, func(files []string) ([]string, error) {
			return FilterFilesByExtension(files, opts.Extension)
		})
	}
	if opts.MaxSize >= 0 {
		filters = append(filters, func(files []string) ([]string, error) {
			return FilterFilesBySize(files, opts.MaxSize)
		})
	}

	return filters
}

func outputFiles(files []string, format string) error {
	switch format {
	case "json":
		jsonData, err := json.Marshal(files)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	default:
		for _, file := range files {
			fmt.Println(file)
		}
	}
	return nil
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
