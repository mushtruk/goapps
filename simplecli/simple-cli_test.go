package simplecli_test

import (
	"bytes"
	"encoding/json"
	"goapps/simplecli"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListFiles(t *testing.T) {
	// Create a new temporary directory
	tempDir, err := os.MkdirTemp("", "testdir_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Create test files in the temporary directory
	_, err = os.Create(filepath.Join(tempDir, "testfile1.txt"))
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_, err = os.Create(filepath.Join(tempDir, "testfile2.txt"))
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Execute the function to test
	files, err := simplecli.ListFiles(tempDir)
	if err != nil {
		t.Fatalf("ListFiles returned an error: %v", err)
	}

	// Sort files to ensure consistent order for comparison
	sort.Strings(files)

	// Verify the results
	expected := []string{"testfile1.txt", "testfile2.txt"}
	sort.Strings(expected) // Sort expected result for consistent comparison
	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Expected %v, got %v", expected, files)
	}
}

func TestFilterFilesByExtension(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir_*")

	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = os.Create(filepath.Join(tempDir, "file1.txt"))
	_, err = os.Create(filepath.Join(tempDir, "file2.jpg"))
	_, err = os.Create(filepath.Join(tempDir, "file3.txt"))

	files, err := simplecli.ListFiles(tempDir)

	txtFiles, err := simplecli.FilterFilesByExtension(files, "txt")

	if err != nil {
		t.Fatalf("FilterFilesByExtension returned an error: %v", err)
	}

	sort.Strings(txtFiles)

	expected := []string{"file1.txt", "file3.txt"}

	if !reflect.DeepEqual(expected, txtFiles) {
		t.Errorf("Expected %v, got %v", expected, txtFiles)
	}
}

func TestFilterFilesBySize(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	createTempFile(t, tempDir, "file1.txt", 100)  // 100 bytes
	createTempFile(t, tempDir, "file2.txt", 2000) // 2000 bytes

	files, err := simplecli.FilterFilesBySize([]string{filepath.Join(tempDir, "file1.txt"), filepath.Join(tempDir, "file2.txt")}, 1500)
	if err != nil {
		t.Fatalf("FilterFilesBySize returned an error: %v", err)
	}

	expected := []string{filepath.Join(tempDir, "file1.txt")}
	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Expected %v, got %v", expected, files)
	}
}

func createTempFile(t *testing.T, dir, name string, size int) {
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if err := f.Truncate(int64(size)); err != nil {
		t.Fatalf("Failed to set size of temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
}

func TestListFilesRecursively(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")

	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	defer os.RemoveAll(tempDir)

	createNestedFiles(t, tempDir, []string{"file1.txt", "subdir/file2.txt", "file2.txt", "subdir1/subdir2/file4.txt"})

	files, err := simplecli.ListFilesRecursively(tempDir)

	expected := []string{"file1.txt", "subdir/file2.txt", "file2.txt", "subdir1/subdir2/file4.txt"}

	sort.Strings(files)
	sort.Strings(expected)

	// Verify
	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Expected %v, got %v", expected, files)
	}
}

func createNestedFiles(t *testing.T, dirName string, paths []string) {
	for _, path := range paths {
		fullPath := filepath.Join(dirName, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)

		if err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}
		_, err = os.Create(fullPath)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}
}

func TestOutputAsJSON(t *testing.T) {
	// Setup: create a temporary directory with some files
	tempDir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	createTempFile(t, tempDir, "file1.txt", 100)
	createTempFile(t, tempDir, "file2.txt", 200)

	// Simulate CLIOptions with JSON output
	opts := simplecli.CLIOptions{
		DirPath: tempDir,
		MaxSize: -1,
		Output:  "json",
		SortMod: false,
	}

	// Capture the standard output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the Init function (or the relevant function that prints the output)
	simplecli.Init(opts)

	// Read the output
	w.Close()
	os.Stdout = originalStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify that the output is correctly formatted as JSON
	var files []string
	err = json.Unmarshal([]byte(output), &files)
	require.NoError(t, err)
	require.Len(t, files, 2)
	require.Contains(t, files, "file1.txt")
	require.Contains(t, files, "file2.txt")
}

func TestSortByModificationTime(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tempDir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create some files with different modification times
	filepath1 := createTempFileWithModTime(t, tempDir, "file1.txt", time.Now())
	filepath2 := createTempFileWithModTime(t, tempDir, "file2.txt", time.Now().Add(-time.Duration(time.Now().Year())))

	testFilePaths := []string{filepath1, filepath2}

	files, err := simplecli.SortFilesByModTime(testFilePaths)
	require.NoError(t, err)

	expected := []string{filepath1, filepath2}
	require.Equal(t, expected, files)
}

func createTempFileWithModTime(t *testing.T, dir, name string, modTime time.Time) string {
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	require.NoError(t, err)
	file.Close()

	err = os.Chtimes(path, modTime, modTime)
	require.NoError(t, err)

	return path
}

func TestFilterFilesByPattern(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tempDir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	createTempFile(t, tempDir, "sample1.txt", 1)
	createTempFile(t, tempDir, "test1.txt", 1)
	createTempFile(t, tempDir, "sample2.txt", 1)
	createTempFile(t, tempDir, "test3.txt", 1)

	opts := simplecli.CLIOptions{
		DirPath: tempDir,
		MaxSize: -1,
		Pattern: "sample.*",
	}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	simplecli.Init(opts)

	w.Close()
	os.Stdout = originalStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Process the text output
	files := []string{}
	if output != "" {
		files = strings.Split(strings.TrimSpace(output), "\n")
	}
	expected := []string{"sample1.txt", "sample2.txt"}
	require.ElementsMatch(t, expected, files)
}
