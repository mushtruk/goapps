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
