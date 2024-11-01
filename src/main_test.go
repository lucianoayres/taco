package main

import (
	"flag"
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
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}
}

// Test for parseArguments function
func TestParseArguments(t *testing.T) {
	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set up initial working directory
	var err error
	initialWorkingDir, err = os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedDirs   []string
		expectedExts   []string
	}{
		{
			name:           "No arguments",
			args:           []string{"taco"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{initialWorkingDir},
			expectedExts:   []string{},
		},
		{
			name:           "Custom output file",
			args:           []string{"taco", "-output=my_output.txt"},
			expectedOutput: "my_output.txt",
			expectedDirs:   []string{initialWorkingDir},
			expectedExts:   []string{},
		},
		{
			name:           "Custom directories",
			args:           []string{"taco", "/path/to/dir1", "/path/to/dir2"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{"/path/to/dir1", "/path/to/dir2"},
			expectedExts:   []string{},
		},
		{
			name:           "Custom output file and directories",
			args:           []string{"taco", "-output=my_output.txt", "/path/to/dir1", "/path/to/dir2"},
			expectedOutput: "my_output.txt",
			expectedDirs:   []string{"/path/to/dir1", "/path/to/dir2"},
			expectedExts:   []string{},
		},
		{
			name:           "Include extensions without leading dots",
			args:           []string{"taco", "-include-ext=go,md"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{initialWorkingDir},
			expectedExts:   []string{".go", ".md"},
		},
		{
			name:           "Include extensions with leading dots",
			args:           []string{"taco", "-include-ext=.go,.md"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{initialWorkingDir},
			expectedExts:   []string{".go", ".md"},
		},
		{
			name:           "Include extensions with mixed case and spaces",
			args:           []string{"taco", "-include-ext=.Go, .MD , .Txt"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{initialWorkingDir},
			expectedExts:   []string{".go", ".md", ".txt"},
		},
		{
			name:           "Include extensions and custom output",
			args:           []string{"taco", "-output=combined.txt", "-include-ext=.js,.json"},
			expectedOutput: "combined.txt",
			expectedDirs:   []string{initialWorkingDir},
			expectedExts:   []string{".js", ".json"},
		},
		{
			name:           "Include extensions and custom directories",
			args:           []string{"taco", "-include-ext=.py,.txt", "/path/to/dir1"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{"/path/to/dir1"},
			expectedExts:   []string{".py", ".txt"},
		},
		{
			name:           "Empty include-ext flag",
			args:           []string{"taco", "-include-ext="},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{initialWorkingDir},
			expectedExts:   []string{},
		},
		{
			name:           "Only include-ext flag with custom output and directories",
			args:           []string{"taco", "-include-ext=.rb,.html", "-output=final.txt", "/dir1", "/dir2"},
			expectedOutput: "final.txt",
			expectedDirs:   []string{"/dir1", "/dir2"},
			expectedExts:   []string{".rb", ".html"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset flag package defaults
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Set os.Args
			os.Args = test.args

			outputFileName, directories, includeExts := parseArguments()
			if outputFileName != test.expectedOutput {
				t.Errorf("Expected output filename %s, got %s", test.expectedOutput, outputFileName)
			}
			if len(directories) != len(test.expectedDirs) {
				t.Errorf("Expected directories %v, got %v", test.expectedDirs, directories)
			} else {
				for i := range directories {
					if directories[i] != test.expectedDirs[i] {
						t.Errorf("Expected directories %v, got %v", test.expectedDirs, directories)
						break
					}
				}
			}
			if len(includeExts) != len(test.expectedExts) {
				t.Errorf("Expected include extensions %v, got %v", test.expectedExts, includeExts)
			} else {
				for i := range includeExts {
					if includeExts[i] != test.expectedExts[i] {
						t.Errorf("Expected include extensions %v, got %v", test.expectedExts, includeExts)
						break
					}
				}
			}
		})
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
		{"another.hidden.file", false}, // Updated to false
		{"nohiddenfile.", false},
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
	tempDir, err := os.MkdirTemp("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	textFile := filepath.Join(tempDir, "textfile.txt")
	binaryFile := filepath.Join(tempDir, "binaryfile.bin")
	emptyFile := filepath.Join(tempDir, "emptyfile.txt")
	nonexistentFile := filepath.Join(tempDir, "nonexistent.txt")

	err = os.WriteFile(textFile, []byte("This is a text file."), 0644)
	if err != nil {
		t.Fatalf("Failed to write text file: %v", err)
	}

	err = os.WriteFile(binaryFile, []byte{0x00, 0xFF, 0x00, 0xFF}, 0644)
	if err != nil {
		t.Fatalf("Failed to write binary file: %v", err)
	}

	err = os.WriteFile(emptyFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}

	if !isTextFile(textFile) {
		t.Errorf("Expected %s to be detected as a text file", textFile)
	}

	if isTextFile(binaryFile) {
		t.Errorf("Expected %s to be detected as a binary file", binaryFile)
	}

	if !isTextFile(emptyFile) {
		t.Errorf("Expected empty file %s to be detected as a text file", emptyFile)
	}

	if isTextFile(nonexistentFile) {
		t.Errorf("Expected nonexistent file %s to be detected as non-text", nonexistentFile)
	}
}

// Test for writeFileContent function
func TestWriteFileContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.txt")
	content := "Hello, Taco!"

	err = os.WriteFile(inputFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = writeFileContent(outFile, inputFile, "input.txt")
	if err != nil {
		t.Errorf("writeFileContent returned error: %v", err)
	}

	outContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "// File: input.txt\n\n" + content + "\n\n"

	if string(outContent) != expectedContent {
		t.Errorf("Output content mismatch.\nGot:\n%s\nExpected:\n%s", string(outContent), expectedContent)
	}
}

func TestWriteFileContent_InputFileError(t *testing.T) {
	// Simulate error when input file does not exist
	tempDir, err := os.MkdirTemp("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	invalidInputFile := filepath.Join(tempDir, "nonexistent.txt")
	outputFile := filepath.Join(tempDir, "output.txt")

	outFile, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = writeFileContent(outFile, invalidInputFile, "nonexistent.txt")
	if err == nil {
		t.Errorf("Expected error when input file does not exist, got nil")
	}
}

func TestWriteFileContent_OutputFileError(t *testing.T) {
	// Simulate error when output file is closed
	tempDir, err := os.MkdirTemp("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.txt")
	content := "Hello, Taco!"

	err = os.WriteFile(inputFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	outFile.Close() // Close the output file to cause an error

	err = writeFileContent(outFile, inputFile, "input.txt")
	if err == nil {
		t.Errorf("Expected error when writing to closed output file, got nil")
	}
}

// Test for processDirectory function
func TestProcessDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	files := map[string]string{
		"file1.txt":                 "Content of file1",
		"file2.go":                  "package main\n\nfunc main() {}",
		".hiddenfile":               "Hidden file content",
		"subdir/file3.md":           "# File3",
		"subdir/.hiddensubfile":     "Hidden subfile content",
		"binaryfile.bin":            string([]byte{0x00, 0xFF}),
		"subdir/binarysubfile.bin":  string([]byte{0x00, 0xFF}),
		"subdir/file4.go":           "package main\n\nfunc helper() {}",
		"subdir/file5.TXT":          "Uppercase extension file",
		"subdir/file6.HTML":         "<html></html>",
		"subdir/file7.py":           "print('Hello, Python')",
		"subdir/file8.NoExtension":  "No extension file",
	}
	createTestFiles(t, tempDir, files)

	outputFilePath := filepath.Join(tempDir, "taco.txt")
	excludedPaths := getExcludedPaths(outputFilePath, "")

	var outputFile *os.File

	// Define includeExts to include .txt and .go files
	includeExts := []string{".txt", ".go"}

	filesProcessed, err := processDirectory(tempDir, &outputFile, outputFilePath, excludedPaths, includeExts)
	if err != nil {
		t.Errorf("processDirectory returned error: %v", err)
	}

	if !filesProcessed {
		t.Errorf("Expected files to be processed")
	}

	// Read the output file
	outputContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check if the content includes expected files
	expectedFiles := []string{
		"file1.txt",
		"file2.go",
		"subdir/file4.go",
		"subdir/file5.TXT", // Should be included because extension is .TXT (lowercased to .txt)
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
		"subdir/file3.md",          // .md not included
		"subdir/file6.HTML",        // .html not included
		"subdir/file7.py",          // .py not included
		"subdir/file8.NoExtension", // No extension not included
	}

	for _, fileName := range unexpectedFiles {
		if strings.Contains(string(outputContent), files[fileName]) {
			t.Errorf("Output content should not include content from %s", fileName)
		}
	}
}

func TestProcessDirectory_ReadDirError(t *testing.T) {
	// Simulate error when reading directory
	invalidDir := "/path/to/nonexistent"
	var outputFile *os.File
	outputFilePath := "output.txt"
	excludedPaths := make(map[string]struct{})

	_, err := processDirectory(invalidDir, &outputFile, outputFilePath, excludedPaths, []string{})
	if err == nil {
		t.Errorf("Expected error when reading nonexistent directory, got nil")
	}
}

// Test for concatenateFiles function
func TestConcatenateFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "taco_test")
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
		"dir1/file2.go":  "package main\n\nfunc main() {}",
		"dir2/file3.md":  "# File3",
		"dir2/file4.txt": "Content of dir2/file4",
		"dir2/file5.js":  "console.log('JavaScript file');",
	}
	createTestFiles(t, tempDir, files)

	outputFilePath := filepath.Join(tempDir, "taco.txt")
	excludedPaths := getExcludedPaths(outputFilePath, "")

	directories := []string{dir1, dir2}

	// Define includeExts to include .txt and .go files
	includeExts := []string{".txt", ".go"}

	err = concatenateFiles(outputFilePath, directories, excludedPaths, includeExts)
	if err != nil {
		t.Errorf("concatenateFiles returned error: %v", err)
	}

	outputContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check if the content includes expected files
	expectedFiles := []string{
		"dir1/file1.txt",
		"dir1/file2.go",
		"dir2/file4.txt",
	}

	for _, filePath := range expectedFiles {
		if !strings.Contains(string(outputContent), files[filePath]) {
			t.Errorf("Output content missing content from %s", filePath)
		}
	}

	// Ensure that .md and .js files are not included
	unexpectedFiles := []string{
		"dir2/file3.md",
		"dir2/file5.js",
	}

	for _, filePath := range unexpectedFiles {
		if strings.Contains(string(outputContent), files[filePath]) {
			t.Errorf("Output content should not include content from %s", filePath)
		}
	}
}

func TestConcatenateFiles_ErrorCases(t *testing.T) {
	// Simulate error when directory does not exist
	tempDir, err := os.MkdirTemp("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	invalidDir := filepath.Join(tempDir, "nonexistent")
	outputFilePath := filepath.Join(tempDir, "taco.txt")
	excludedPaths := getExcludedPaths(outputFilePath, "")

	err = concatenateFiles(outputFilePath, []string{invalidDir}, excludedPaths, []string{})
	if err == nil {
		t.Errorf("Expected error for nonexistent directory, got nil")
	}
}

func TestConcatenateFiles_OutputFileError(t *testing.T) {
	// Simulate error when output file cannot be created
	tempDir, err := os.MkdirTemp("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a text file in tempDir to ensure the output file is attempted to be created
	testFilePath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a directory where the output file should be, causing an error when trying to write
	outputFilePath := filepath.Join(tempDir, "taco.txt")
	err = os.Mkdir(outputFilePath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory %s: %v", outputFilePath, err)
	}

	excludedPaths := getExcludedPaths(outputFilePath, "")

	err = concatenateFiles(outputFilePath, []string{tempDir}, excludedPaths, []string{})
	if err == nil {
		t.Errorf("Expected error when output file path is a directory, got nil")
	}

	// Defer the removal of the output directory
	defer func() {
		if err := os.RemoveAll(outputFilePath); err != nil && !os.IsNotExist(err) {
			t.Errorf("Failed to remove output directory: %v", err)
		}
	}()
}

// Test for run function (refactored main logic)
func TestRun(t *testing.T) {
	// Save original os.Args, flag.CommandLine, and working directory
	originalArgs := os.Args
	originalFlagCommandLine := flag.CommandLine
	originalWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlagCommandLine
		os.Chdir(originalWorkingDir)
	}()

	tempDir, err := os.MkdirTemp("", "taco_test_run")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	files := map[string]string{
		"file1.txt": "Content of file1",
		"file2.go":  "package main\n\nfunc main() {}",
		"file3.md":  "# File3",
	}
	createTestFiles(t, tempDir, files)

	// Set os.Args to simulate command-line arguments with -include-ext
	os.Args = []string{"taco", "-output=taco_output.txt", "-include-ext=.txt,.go", tempDir}

	// Reset flag package defaults
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Change the working directory to tempDir
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Run the main logic
	err = run()
	if err != nil {
		t.Errorf("run returned error: %v", err)
	}

	// Verify that output file is created with expected content
	outputFilePath := filepath.Join(tempDir, "taco_output.txt")
	outputContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check if the content matches expected content
	expectedFiles := []string{"file1.txt", "file2.go"}
	for _, fileName := range expectedFiles {
		if !strings.Contains(string(outputContent), files[fileName]) {
			t.Errorf("Output content missing content from %s", fileName)
		}
	}

	// Ensure that .md files are not included
	if strings.Contains(string(outputContent), files["file3.md"]) {
		t.Errorf("Output content should not include content from %s", "file3.md")
	}
}

// Test for getExcludedPaths function
func TestGetExcludedPaths(t *testing.T) {
	outputFilePath := "/path/to/output.txt"
	scriptFilePath := "/path/to/script.go"
	expected := map[string]struct{}{
		outputFilePath: {},
		scriptFilePath: {},
	}

	result := getExcludedPaths(outputFilePath, scriptFilePath)
	if len(result) != len(expected) {
		t.Errorf("Expected excluded paths %v, got %v", expected, result)
	}
	for key := range expected {
		if _, exists := result[key]; !exists {
			t.Errorf("Expected path %s to be excluded", key)
		}
	}
}

// Additional Tests for include-ext functionality

// TestShouldIncludeFile tests the shouldIncludeFile helper function
func TestShouldIncludeFile(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		includeExts []string
		expected    bool
	}{
		{
			name:        "No includeExts, should include all",
			filePath:    "file1.txt",
			includeExts: []string{},
			expected:    true,
		},
		{
			name:        "Include .txt, file with .txt extension",
			filePath:    "file1.txt",
			includeExts: []string{".txt"},
			expected:    true,
		},
		{
			name:        "Include .txt, file with .go extension",
			filePath:    "file2.go",
			includeExts: []string{".txt"},
			expected:    false,
		},
		{
			name:        "Include .go and .md, file with .md extension",
			filePath:    "file3.md",
			includeExts: []string{".go", ".md"},
			expected:    true,
		},
		{
			name:        "Include .go and .md (lowercased), file with .go extension",
			filePath:    "file2.go",
			includeExts: []string{".go", ".md"},
			expected:    true,
		},
		{
			name:        "Include .txt and .go, file with uppercase extension",
			filePath:    "file1.TXT",
			includeExts: []string{".txt", ".go"},
			expected:    true,
		},
		{
			name:        "Include .txt and .go, file with .TXT extension",
			filePath:    "file1.TXT",
			includeExts: []string{".txt", ".go"},
			expected:    true,
		},
		{
			name:        "Include .txt and .go, file with .Md extension",
			filePath:    "file3.Md",
			includeExts: []string{".txt", ".go"},
			expected:    false,
		},
		{
			name:        "Include .txt and .go, file with no extension",
			filePath:    "file4",
			includeExts: []string{".txt", ".go"},
			expected:    false,
		},
	}

	for _, test := range tests {
		result := shouldIncludeFile(test.filePath, test.includeExts)
		if result != test.expected {
			t.Errorf("shouldIncludeFile(%q, %v) = %v; want %v", test.filePath, test.includeExts, result, test.expected)
		}
	}
}

// Test shouldIncludeFile with various edge cases
func TestShouldIncludeFile_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		includeExts []string
		expected    bool
	}{
		{
			name:        "Empty filePath",
			filePath:    "",
			includeExts: []string{".txt"},
			expected:    false,
		},
		{
			name:        "File with multiple dots",
			filePath:    "archive.tar.gz",
			includeExts: []string{".gz"},
			expected:    true,
		},
		{
			name:        "File with uppercase extension included",
			filePath:    "README.MD",
			includeExts: []string{".md"},
			expected:    true, // Because includeExts are already lowercased
		},
		{
			name:        "File with extension not in includeExts",
			filePath:    "script.py",
			includeExts: []string{".js", ".ts"},
			expected:    false,
		},
	}

	for _, test := range tests {
		result := shouldIncludeFile(test.filePath, test.includeExts)
		if result != test.expected {
			t.Errorf("shouldIncludeFile(%q, %v) = %v; want %v", test.filePath, test.includeExts, result, test.expected)
		}
	}
}
