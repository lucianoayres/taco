package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseArguments checks if parseArguments correctly parses command-line arguments.
func TestParseArguments(t *testing.T) {
	// Define flags for testing
	os.Args = []string{"cmd", "-output", "test_output.txt", "-include-ext", ".go,.md", "-exclude-ext", ".test", "-exclude-dir", "vendor"}
	
	// Reset flag defaults and parse
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	output, _, includeExt, excludeExt, excludeDir, verbose, err := parseArguments() // Remove unused 'directories' variable
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if output != "test_output.txt" {
		t.Errorf("Expected output 'test_output.txt', got %s", output)
	}
	if !strings.Contains(includeExt[0], ".go") || !strings.Contains(includeExt[1], ".md") {
		t.Errorf("Expected include extensions [.go .md], got %v", includeExt)
	}
	if excludeExt[0] != ".test" {
		t.Errorf("Expected exclude extension '.test', got %v", excludeExt)
	}
	if excludeDir[0] != "vendor" {
		t.Errorf("Expected exclude directory 'vendor', got %v", excludeDir)
	}
	if verbose {
		t.Errorf("Expected verbose to be false, got true")
	}
}

// TestShouldIncludeFile checks the file inclusion logic based on file extensions.
func TestShouldIncludeFile(t *testing.T) {
	includeExts := []string{".go", ".md"}
	excludeExts := []string{".test", ".spec"}
	
	if !shouldIncludeFile("main.go", includeExts, excludeExts) {
		t.Error("Expected 'main.go' to be included")
	}
	if shouldIncludeFile("test.spec", includeExts, excludeExts) {
		t.Error("Expected 'test.spec' to be excluded")
	}
}

// TestGetExcludedPaths ensures script and output paths are excluded from concatenation.
func TestGetExcludedPaths(t *testing.T) {
	excludePaths := getExcludedPaths("output.txt", "script.go")
	if _, ok := excludePaths["output.txt"]; !ok {
		t.Error("Expected output.txt to be excluded")
	}
	if _, ok := excludePaths["script.go"]; !ok {
		t.Error("Expected script.go to be excluded")
	}
}

// TestIsTextFile verifies text file detection using MIME type.
func TestIsTextFile(t *testing.T) {
	// Create temporary file for testing
	file, err := os.CreateTemp("", "test.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())
	file.WriteString("This is a test text file.")

	if !isTextFile(file.Name()) {
		t.Error("Expected text file detection to be true for test.txt")
	}
}

// TestWriteFileContent ensures content is written with file path annotations.
func TestWriteFileContent(t *testing.T) {
	// Set up test files and output
	contentFile, _ := os.CreateTemp("", "content.txt")
	defer os.Remove(contentFile.Name())
	contentFile.WriteString("File content")

	outputFile, _ := os.CreateTemp("", "output.txt")
	defer os.Remove(outputFile.Name())

	// Run writeFileContent
	err := writeFileContent(outputFile, contentFile.Name(), "content.txt")
	if err != nil {
		t.Fatalf("Error writing file content: %v", err)
	}

	// Verify content format in output
	data, _ := os.ReadFile(outputFile.Name())
	expected := "// File: content.txt\n\nFile content\n"
	if string(data) != expected {
		t.Errorf("Expected output format:\n%s\nGot:\n%s", expected, data)
	}
}

// TestConcatenateFiles validates concatenation of files in directories.
func TestConcatenateFiles(t *testing.T) {
	// Create directories and files for testing
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.go")
	file2 := filepath.Join(dir, "file2.md")
	outputFile := filepath.Join(dir, "output.txt")
	os.WriteFile(file1, []byte("Content of file1"), 0644)
	os.WriteFile(file2, []byte("Content of file2"), 0644)

	// Exclude paths map
	excludedPaths := map[string]struct{}{outputFile: {}, "exclude_dir": {}}

	// Run concatenateFiles
	err := concatenateFiles(outputFile, []string{dir}, excludedPaths, nil, []string{".go", ".md"}, nil, true)
	if err != nil {
		t.Fatalf("Error concatenating files: %v", err)
	}

	// Check if output file contains both files
	data, _ := os.ReadFile(outputFile)

	// Adjust expected output to include the absolute paths of file1 and file2
	expected := "// File: " + file1 + "\n\nContent of file1\n// File: " + file2 + "\n\nContent of file2\n"
	if string(data) != expected {
		t.Errorf("Expected concatenated output:\n%s\nGot:\n%s", expected, data)
	}
}
