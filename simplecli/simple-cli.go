package simplecli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type CLIOptions struct {
	DirPath   string
	Extension string
	MaxSize   int
}

func ParseFlags() CLIOptions {
	dirPath := flag.String("path", ".", "Directory path to list files")
	extension := flag.String("ext", "", "File extension to filter by (e.g., .txt)")
	size := flag.Int("size", -1, "File size to filter by")

	flag.Parse()

	return CLIOptions{
		DirPath:   *dirPath,
		Extension: *extension,
		MaxSize:   *size,
	}
}

func Init(opts CLIOptions) {
	files, err := ListFiles(opts.DirPath)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	if opts.Extension != "" {
		files, err = FilterFilesByExtension(files, opts.Extension)
		if err != nil {
			log.Fatalf("Error filtering files: %v", err)
		}
	}

	if opts.MaxSize >= 0 {
		files, err = FilterFilesBySize(files, opts.MaxSize)
		if err != nil {
			log.Fatalf("Error filtering files: %v", err)
		}
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
