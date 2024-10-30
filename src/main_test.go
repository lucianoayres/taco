package main

import (
	"flag"
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
	}{
		{
			name:           "No arguments",
			args:           []string{"taco"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{initialWorkingDir},
		},
		{
			name:           "Custom output file",
			args:           []string{"taco", "-output=my_output.txt"},
			expectedOutput: "my_output.txt",
			expectedDirs:   []string{initialWorkingDir},
		},
		{
			name:           "Custom directories",
			args:           []string{"taco", "/path/to/dir1", "/path/to/dir2"},
			expectedOutput: "taco.txt",
			expectedDirs:   []string{"/path/to/dir1", "/path/to/dir2"},
		},
		{
			name:           "Custom output file and directories",
			args:           []string{"taco", "-output=my_output.txt", "/path/to/dir1", "/path/to/dir2"},
			expectedOutput: "my_output.txt",
			expectedDirs:   []string{"/path/to/dir1", "/path/to/dir2"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset flag package defaults
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Set os.Args
			os.Args = test.args

			outputFileName, directories := parseArguments()
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

	ioutil.WriteFile(textFile, []byte("This is a text file."), 0644)
	ioutil.WriteFile(binaryFile, []byte{0x00, 0xFF, 0x00, 0xFF}, 0644)
	ioutil.WriteFile(emptyFile, []byte(""), 0644)

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
	tempDir, err := ioutil.TempDir("", "taco_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.txt")
	content := "Hello, Taco!"

	os.WriteFile(inputFile, []byte(content), 0644)

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
	outputContent, err := os.ReadFile(outputFilePath)
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

func TestProcessDirectory_ReadDirError(t *testing.T) {
	// Simulate error when reading directory
	invalidDir := "/path/to/nonexistent"
	var outputFile *os.File
	outputFilePath := "output.txt"
	excludedPaths := make(map[string]struct{})

	_, err := processDirectory(invalidDir, &outputFile, outputFilePath, excludedPaths)
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

	outputContent, err := os.ReadFile(outputFilePath)
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

	err = concatenateFiles(outputFilePath, []string{invalidDir}, excludedPaths)
	if err == nil {
		t.Errorf("Expected error for nonexistent directory, got nil")
	}
}

func TestConcatenateFiles_OutputFileError(t *testing.T) {
    // Simulate error when output file cannot be created
    tempDir, err := ioutil.TempDir("", "taco_test")
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
    if err := os.Mkdir(outputFilePath, 0755); err != nil {
        t.Fatalf("Failed to create directory %s: %v", outputFilePath, err)
    }

    excludedPaths := getExcludedPaths(outputFilePath, "")

    err = concatenateFiles(outputFilePath, []string{tempDir}, excludedPaths)
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
        "file2.txt": "Content of file2",
    }
    createTestFiles(t, tempDir, files)

    // Set os.Args to simulate command-line arguments
    os.Args = []string{"taco", "-output=taco_output.txt", tempDir}

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
    for _, fileName := range []string{"file1.txt", "file2.txt"} {
        if !strings.Contains(string(outputContent), files[fileName]) {
            t.Errorf("Output content missing content from %s", fileName)
        }
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
