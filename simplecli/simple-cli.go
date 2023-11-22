package simplecli

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type CLIOptions struct {
	DirPath   string
	Extension string
	MaxSize   int
	Recursive bool
	SortMod   bool
	Output    string
	Pattern   string
	Checksum  string
}

func ParseFlags() CLIOptions {
	dirPath := flag.String("path", ".", "Directory path to list files")
	extension := flag.String("ext", "", "File extension to filter by (e.g., .txt)")
	size := flag.Int("size", -1, "File size to filter by")
	recursive := flag.Bool("recursive", false, "List files in directories recursively")
	output := flag.String("output", "text", "Output format (e.g., 'text', 'JSON'")
	sortByModTime := flag.Bool("sortByModTime", false, "Sort files by modification time")
	pattern := flag.String("pattern", "", "Regex pattern to filter files")
	checksum := flag.String("checksum", "", "Calculate the MD5 checksum of a specified file")

	flag.Parse()

	return CLIOptions{
		DirPath:   *dirPath,
		Extension: *extension,
		MaxSize:   *size,
		Recursive: *recursive,
		SortMod:   *sortByModTime,
		Output:    *output,
		Pattern:   *pattern,
		Checksum:  *checksum,
	}
}

type FilterFunc func([]string) ([]string, error)

func Init(opts CLIOptions) {
	files, err := listFilesBasedOnFlag(opts)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	if opts.Checksum != "" {
		checksum, err := CalculateMD5Checksum(opts.Checksum)
		if err != nil {
			log.Fatalf("Error calculating checksum: %v", err)
		}
		fmt.Println("MD5 Checksum:", checksum)
		return
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

	if opts.SortMod {
		filters = append(filters, func(files []string) ([]string, error) {
			return SortFilesByModTime(files)
		})
	}

	if opts.Pattern != "" {
		filters = append(filters, func(files []string) ([]string, error) {
			return FilterFilesByPattern(files, opts.Pattern)
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
		var sb strings.Builder
		for _, file := range files {
			sb.WriteString(file + "\n")
		}
		fmt.Print(sb.String())
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

func SortFilesByModTime(filePaths []string) ([]string, error) {
	type fileModTime struct {
		path    string
		modTime time.Time
	}

	var filesWithModTime []fileModTime

	for _, file := range filePaths {
		info, err := os.Stat(file)

		if err != nil {
			return nil, err
		}

		filesWithModTime = append(filesWithModTime, fileModTime{
			path:    file,
			modTime: info.ModTime(),
		})
	}

	sort.Slice(filesWithModTime, func(i, j int) bool {
		return filesWithModTime[i].modTime.Before(filesWithModTime[i].modTime)
	})

	var sortedFiles []string

	for _, file := range filesWithModTime {
		sortedFiles = append(sortedFiles, file.path)
	}

	return sortedFiles, nil
}

func FilterFilesByPattern(files []string, pattern string) ([]string, error) {
	var filteredFiles []string
	regex, err := regexp.Compile(pattern)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if regex.MatchString(file) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, err
}

func CalculateMD5Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hasher := md5.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
