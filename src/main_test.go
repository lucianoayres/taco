package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// TestParseArguments verifies that arguments are parsed correctly
func TestParseArguments(t *testing.T) {
	// Set up command line arguments
	os.Args = []string{"cmd", "-output", "output.txt", "-include-ext", ".go,.txt", "-exclude-ext", ".md"}
	output, directories, include, exclude, _ := parseArguments()

	// Check the parsed arguments
	if output != "output.txt" {
		t.Errorf("Expected output filename 'output.txt', got %s", output)
	}
	if len(directories) == 0 || directories[0] != initialWorkingDir {
		t.Errorf("Expected default directory to be %s, got %v", initialWorkingDir, directories)
	}
	expectedInclude := []string{".go", ".txt"}
	if !equalSlices(include, expectedInclude) {
		t.Errorf("Expected include extensions %v, got %v", expectedInclude, include)
	}
	expectedExclude := []string{".md"}
	if !equalSlices(exclude, expectedExclude) {
		t.Errorf("Expected exclude extensions %v, got %v", expectedExclude, exclude)
	}
}

// TestGetExcludedPaths checks if the script and output paths are excluded
func TestGetExcludedPaths(t *testing.T) {
	outputFile := "output.txt"
	scriptFile := "script.go"
	excludedPaths := getExcludedPaths(outputFile, scriptFile)
	if len(excludedPaths) != 2 {
		t.Errorf("Expected 2 excluded paths, got %d", len(excludedPaths))
	}
	if _, exists := excludedPaths[outputFile]; !exists {
		t.Errorf("Output file %s should be excluded", outputFile)
	}
	if _, exists := excludedPaths[scriptFile]; !exists {
		t.Errorf("Script file %s should be excluded", scriptFile)
	}
}

// TestConcatenateFiles verifies that files are concatenated properly
func TestConcatenateFiles(t *testing.T) {
	// Create a temporary directory and files
	tempDir := t.TempDir()
	originalWorkingDir := initialWorkingDir
	initialWorkingDir = tempDir // Set initialWorkingDir to tempDir for consistency
	defer func() {
		initialWorkingDir = originalWorkingDir
	}()

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")

	// Write to files
	err := os.WriteFile(file1, []byte("Content of file1"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to file1: %v", err)
	}
	err = os.WriteFile(file2, []byte("Content of file2"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to file2: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.txt")
	excludedPaths := getExcludedPaths(outputFile, "")

	// Test concatenation
	err = concatenateFiles(outputFile, []string{tempDir}, excludedPaths, map[string]struct{}{}, []string{}, []string{})
	if err != nil {
		t.Fatalf("Concatenation failed: %v", err)
	}

	// Read output file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Compute relative paths for expected content
	relFile1, err := filepath.Rel(initialWorkingDir, file1)
	if err != nil {
		relFile1 = file1
	}
	relFile2, err := filepath.Rel(initialWorkingDir, file2)
	if err != nil {
		relFile2 = file2
	}

	// Adjust expected content to use relative paths and a single newline separator
	expectedContent := fmt.Sprintf("// File: %s\n\nContent of file1\n// File: %s\n\nContent of file2\n", relFile1, relFile2)
	if string(data) != expectedContent {
		t.Errorf("Expected concatenated content:\n%s\nGot:\n%s", expectedContent, string(data))
	}
}

// TestProcessDirectory verifies processing of a single directory
func TestProcessDirectory(t *testing.T) {
	tempDir := t.TempDir()
	originalWorkingDir := initialWorkingDir
	initialWorkingDir = tempDir // Set initialWorkingDir to tempDir for consistency
	defer func() {
		initialWorkingDir = originalWorkingDir
	}()

	file := filepath.Join(tempDir, "file.txt")
	err := os.WriteFile(file, []byte("Content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.txt")
	excludedPaths := getExcludedPaths(outputFile, "")

	var outputFileHandler *os.File
	defer func() {
		if outputFileHandler != nil {
			outputFileHandler.Close()
		}
	}()

	filesProcessed, err := processDirectory(tempDir, &outputFileHandler, outputFile, excludedPaths, map[string]struct{}{}, []string{}, []string{})
	if err != nil {
		t.Fatalf("Error processing directory: %v", err)
	}
	if !filesProcessed {
		t.Error("Expected files to be processed")
	}

	// Verify the content of the output file
	if outputFileHandler != nil {
		data, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		relFile, err := filepath.Rel(initialWorkingDir, file)
		if err != nil {
			relFile = file
		}
		expectedContent := fmt.Sprintf("// File: %s\n\nContent\n", relFile)
		if string(data) != expectedContent {
			t.Errorf("Expected output content:\n%s\nGot:\n%s", expectedContent, string(data))
		}
	}
}

// TestShouldIncludeFile tests whether files should be included based on extensions
func TestShouldIncludeFile(t *testing.T) {
	includeExts := []string{".txt"}
	excludeExts := []string{".md"}
	tests := []struct {
		filePath string
		expected bool
	}{
		{"file.txt", true},
		{"file.md", false},
		{"file.go", false},
		{"file.txt.md", false},
	}

	for _, test := range tests {
		result := shouldIncludeFile(test.filePath, includeExts, excludeExts)
		if result != test.expected {
			t.Errorf("Expected shouldIncludeFile(%s) to return %v, got %v", test.filePath, test.expected, result)
		}
	}
}

// TestIsTextFile tests if a file is detected as a text file
func TestIsTextFile(t *testing.T) {
	tempDir := t.TempDir()
	originalWorkingDir := initialWorkingDir
	initialWorkingDir = tempDir // Set initialWorkingDir to tempDir for consistency
	defer func() {
		initialWorkingDir = originalWorkingDir
	}()

	// Text file
	textFile := filepath.Join(tempDir, "text.txt")
	err := os.WriteFile(textFile, []byte("This is a text file"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to text file: %v", err)
	}
	if !isTextFile(textFile) {
		t.Error("Expected text file to be detected as text")
	}

	// Non-text file
	binaryFile := filepath.Join(tempDir, "binary.bin")
	err = os.WriteFile(binaryFile, []byte{0xFF, 0xD8, 0xFF}, 0644) // JPEG header
	if err != nil {
		t.Fatalf("Failed to write to binary file: %v", err)
	}
	if isTextFile(binaryFile) {
		t.Error("Expected binary file to be detected as non-text")
	}
}

// TestWriteFileContent verifies the writing of a file's content
func TestWriteFileContent(t *testing.T) {
	tempDir := t.TempDir()
	originalWorkingDir := initialWorkingDir
	initialWorkingDir = tempDir // Set initialWorkingDir to tempDir for consistency
	defer func() {
		initialWorkingDir = originalWorkingDir
	}()

	file := filepath.Join(tempDir, "file.txt")
	outputFile := filepath.Join(tempDir, "output.txt")
	err := os.WriteFile(file, []byte("File content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	output, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Error creating output file: %v", err)
	}
	defer output.Close()

	err = writeFileContent(output, file, "file.txt")
	if err != nil {
		t.Fatalf("Error writing file content: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Error reading output file: %v", err)
	}
	expected := "// File: file.txt\n\nFile content\n"
	if string(data) != expected {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expected, string(data))
	}
}

// Helper function to compare two string slices
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
