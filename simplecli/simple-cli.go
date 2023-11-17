package simplecli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Command interface {
	Execute() ([]string, error)
}

type ListFilesCommand struct {
	DirPath string
}

func (c ListFilesCommand) Execute() ([]string, error) {
	return ListFiles(c.DirPath)
}

type FilterFilesByExtensionCommand struct {
	DirPath   string
	Extension string
}

type CLIOptions struct {
	DirPath   string
	Extension string
}

func ParseFlags() CLIOptions {
	dirPath := flag.String("path", ".", "Directory path to list files")
	extension := flag.String("ext", "", "File extension to filter by (e.g., .txt)")

	flag.Parse()

	return CLIOptions{
		DirPath:   *dirPath,
		Extension: *extension,
	}
}

func Init(opts CLIOptions) {
	files, err := ListFiles(opts.DirPath)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	if opts.Extension != "" {
		files, err = FilterFiles(files, opts.Extension)
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

func FilterFiles(files []string, extension string) ([]string, error) {
	var filteredFiles []string
	for _, file := range files {
		if strings.HasSuffix(file, extension) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}
