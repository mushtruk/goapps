package simplecli_test

import (
	"goapps/simplecli"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
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
