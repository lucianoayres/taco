package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create temporary files and directories for testing
func createTestFiles(t *testing.T, baseDir string, files map[string]string) {
	for path, content := range files {
		fullPath := filepath.Join(baseDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}
}

// Test for isHidden function
func TestIsHidden(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{".hiddenfile", true},
		{"visiblefile", false},
		{".hiddenDir/", true},
		{"visibleDir/", false},
	}

	for _, test := range tests {
		result := isHidden(test.name)
		if result != test.expected {
			t.Errorf("isHidden(%q) = %v; want %v", test.name, result, test.expected)
		}
	}
}

// Test for isTextFile function
func TestIsTextFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	textFile := filepath.Join(tempDir, "textfile.txt")
	binaryFile := filepath.Join(tempDir, "binaryfile.bin")

	ioutil.WriteFile(textFile, []byte("This is a text file."), 0644)
	ioutil.WriteFile(binaryFile, []byte{0x00, 0xFF, 0x00, 0xFF}, 0644)

	if !isTextFile(textFile) {
		t.Errorf("Expected %s to be detected as a text file", textFile)
	}

	if isTextFile(binaryFile) {
		t.Errorf("Expected %s to be detected as a binary file", binaryFile)
	}
}

// Test for writeFileContent function
func TestWriteFileContent(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.txt")
	content := "Hello, Taco!"

	ioutil.WriteFile(inputFile, []byte(content), 0644)

	outFile, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = writeFileContent(outFile, inputFile, "input.txt")
	if err != nil {
		t.Errorf("writeFileContent returned error: %v", err)
	}

	outContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "// File: input.txt\n\n" + content + "\n\n"

	if string(outContent) != expectedContent {
		t.Errorf("Output content mismatch.\nGot:\n%s\nExpected:\n%s", string(outContent), expectedContent)
	}
}

// Test for processDirectory function
func TestProcessDirectory(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	files := map[string]string{
		"file1.txt":                "Content of file1",
		"file2.txt":                "Content of file2",
		".hiddenfile":              "Hidden file content",
		"subdir/file3.txt":         "Content of file3",
		"subdir/.hiddensubfile":    "Hidden subfile content",
		"binaryfile.bin":           string([]byte{0x00, 0xFF}),
		"subdir/binarysubfile.bin": string([]byte{0x00, 0xFF}),
	}
	createTestFiles(t, tempDir, files)

	outputFilePath := filepath.Join(tempDir, "taco.txt")
	excludedPaths := getExcludedPaths(outputFilePath, "")

	var outputFile *os.File

	filesProcessed, err := processDirectory(tempDir, &outputFile, outputFilePath, excludedPaths)
	if err != nil {
		t.Errorf("processDirectory returned error: %v", err)
	}

	if !filesProcessed {
		t.Errorf("Expected files to be processed")
	}

	// Read the output file
	outputContent, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check if the content matches expected content
	expectedFiles := []string{
		"file1.txt",
		"file2.txt",
		"subdir/file3.txt",
	}
	for _, fileName := range expectedFiles {
		if !strings.Contains(string(outputContent), files[fileName]) {
			t.Errorf("Output content missing content from %s", fileName)
		}
	}

	// Ensure hidden and binary files are not included
	unexpectedFiles := []string{
		".hiddenfile",
		"subdir/.hiddensubfile",
		"binaryfile.bin",
		"subdir/binarysubfile.bin",
	}
	for _, fileName := range unexpectedFiles {
		if strings.Contains(string(outputContent), files[fileName]) {
			t.Errorf("Output content should not include content from %s", fileName)
		}
	}
}

// Test for concatenateFiles function
func TestConcatenateFiles(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files in multiple directories
	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir2")
	os.Mkdir(dir1, 0755)
	os.Mkdir(dir2, 0755)

	files := map[string]string{
		"dir1/file1.txt": "Content of dir1/file1",
		"dir2/file2.txt": "Content of dir2/file2",
	}
	createTestFiles(t, tempDir, files)

	outputFilePath := filepath.Join(tempDir, "taco.txt")
	excludedPaths := getExcludedPaths(outputFilePath, "")

	directories := []string{dir1, dir2}

	err = concatenateFiles(outputFilePath, directories, excludedPaths)
	if err != nil {
		t.Errorf("concatenateFiles returned error: %v", err)
	}

	outputContent, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check if the content matches expected content
	for _, filePath := range []string{"dir1/file1.txt", "dir2/file2.txt"} {
		if !strings.Contains(string(outputContent), files[filePath]) {
			t.Errorf("Output content missing content from %s", filePath)
		}
	}
}