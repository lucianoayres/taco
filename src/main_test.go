// File: src/main_test.go

package main

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestParseArguments checks if parseArguments correctly parses command-line arguments.
func TestParseArguments(t *testing.T) {
	// Define flags for testing
	os.Args = []string{
		"cmd",
		"-output", "test_output.txt",
		"-include-ext", ".go,.md",
		"-exclude-ext", ".test",
		"-include-file-pattern", "^main\\.go$",
		"-exclude-dir", "vendor",
		"-exclude-file-pattern", ".*_test\\.go$",
		"-verbose",
	}

	// Reset flag defaults and parse
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	output, _, includeExt, excludeExt, includePatterns, excludePatterns, excludeDir, verbose, err := parseArguments()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if output != "test_output.txt" {
		t.Errorf("Expected output 'test_output.txt', got %s", output)
	}
	if len(includeExt) != 2 || includeExt[0] != ".go" || includeExt[1] != ".md" {
		t.Errorf("Expected include extensions [.go .md], got %v", includeExt)
	}
	if len(excludeExt) != 1 || excludeExt[0] != ".test" {
		t.Errorf("Expected exclude extension '.test', got %v", excludeExt)
	}
	if len(excludeDir) != 1 || excludeDir[0] != "vendor" {
		t.Errorf("Expected exclude directory 'vendor', got %v", excludeDir)
	}
	if len(includePatterns) != 1 || includePatterns[0] != "^main\\.go$" {
		t.Errorf("Expected include pattern '^main\\.go$', got %v", includePatterns)
	}
	if len(excludePatterns) != 1 || excludePatterns[0] != ".*_test\\.go$" {
		t.Errorf("Expected exclude pattern '.*_test\\.go$', got %v", excludePatterns)
	}
	if !verbose {
		t.Errorf("Expected verbose to be true, got false")
	}
}

// TestShouldIncludeFile checks the file inclusion logic based on file extensions.
func TestShouldIncludeFile(t *testing.T) {
	includeExts := []string{".go", ".md"}
	excludeExts := []string{".test", ".spec"}

	if !shouldIncludeFile("main.go", includeExts, excludeExts) {
		t.Error("Expected 'main.go' to be included")
	}
	if !shouldIncludeFile("README.md", includeExts, excludeExts) {
		t.Error("Expected 'README.md' to be included")
	}
	if shouldIncludeFile("test.spec", includeExts, excludeExts) {
		t.Error("Expected 'test.spec' to be excluded")
	}
	if shouldIncludeFile("example.test", includeExts, excludeExts) {
		t.Error("Expected 'example.test' to be excluded")
	}
}

// TestMatchesPatterns checks if filenames are correctly matched against include and exclude patterns.
func TestMatchesPatterns(t *testing.T) {
	includePatterns := []string{"^main\\.go$", "^README\\.md$"}
	excludePatterns := []string{".*_test\\.go$", "^LICENSE$"}

	var includeRegexps []*regexp.Regexp
	for _, pattern := range includePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			t.Fatalf("Error compiling include pattern %q: %v", pattern, err)
		}
		includeRegexps = append(includeRegexps, re)
	}

	var excludeRegexps []*regexp.Regexp
	for _, pattern := range excludePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			t.Fatalf("Error compiling exclude pattern %q: %v", pattern, err)
		}
		excludeRegexps = append(excludeRegexps, re)
	}

	tests := []struct {
		filename      string
		shouldInclude bool
	}{
		{"main.go", true},
		{"README.md", true},
		{"main_test.go", false},
		{"LICENSE", false},
		{"utils.go", false},
		{"doc.md", false},
	}

	for _, test := range tests {
		includeMatch := matchesPatterns(test.filename, includeRegexps)
		excludeMatch := matchesPatterns(test.filename, excludeRegexps)
		shouldInclude := (len(includeRegexps) == 0 || includeMatch) && !excludeMatch

		if shouldInclude != test.shouldInclude {
			t.Errorf("Expected file %q inclusion to be %v, got %v", test.filename, test.shouldInclude, shouldInclude)
		}
	}
}

// TestIsTextFile verifies text file detection using MIME type.
func TestIsTextFile(t *testing.T) {
	// Create temporary text file
	textFile, err := os.CreateTemp("", "test.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(textFile.Name())
	textFile.WriteString("This is a test text file.")

	if !isTextFile(textFile.Name()) {
		t.Error("Expected text file detection to be true for test.txt")
	}

	// Create temporary binary file
	binaryFile, err := os.CreateTemp("", "test.bin")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(binaryFile.Name())
	binaryFile.Write([]byte{0x00, 0xFF, 0x00, 0xFF})

	if isTextFile(binaryFile.Name()) {
		t.Error("Expected text file detection to be false for test.bin")
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
	// Create a temporary parent directory
	parentDir := t.TempDir()
	// Create a subdirectory in the parent directory
	dirName := "testdir"
	dir := filepath.Join(parentDir, dirName)
	os.Mkdir(dir, 0755)

	// Create files in the subdirectory
	file1 := filepath.Join(dir, "main.go")
	file2 := filepath.Join(dir, "README.md")
	file3 := filepath.Join(dir, "utils.go")
	file4 := filepath.Join(dir, "main_test.go")
	outputFile := filepath.Join(parentDir, "output.txt")
	os.WriteFile(file1, []byte("Content of main.go"), 0644)
	os.WriteFile(file2, []byte("Content of README.md"), 0644)
	os.WriteFile(file3, []byte("Content of utils.go"), 0644)
	os.WriteFile(file4, []byte("Content of main_test.go"), 0644)

	// Set initialWorkingDir to the parent directory
	initialWorkingDir = parentDir

	// Exclude paths map
	excludedPaths := map[string]struct{}{outputFile: {}}

	// Include and exclude patterns
	includePatterns := []string{"^main\\.go$", "^README\\.md$"}
	excludePatterns := []string{".*_test\\.go$"}

	var includeRegexps []*regexp.Regexp
	for _, pattern := range includePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			t.Fatalf("Error compiling include pattern %q: %v", pattern, err)
		}
		includeRegexps = append(includeRegexps, re)
	}

	var excludeRegexps []*regexp.Regexp
	for _, pattern := range excludePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			t.Fatalf("Error compiling exclude pattern %q: %v", pattern, err)
		}
		excludeRegexps = append(excludeRegexps, re)
	}

	// Run concatenateFiles with the relative directory name
	err := concatenateFiles(outputFile, []string{dirName}, excludedPaths, nil, []string{".go", ".md"}, nil, includeRegexps, excludeRegexps, true)
	if err != nil {
		t.Fatalf("Error concatenating files: %v", err)
	}

	// Check if output file contains only file1 and file2
	data, _ := os.ReadFile(outputFile)

	// Expected output (files in alphabetical order)
	expected := "// File: testdir/README.md\n\nContent of README.md\n// File: testdir/main.go\n\nContent of main.go\n"
	if string(data) != expected {
		t.Errorf("Expected concatenated output:\n%s\nGot:\n%s", expected, data)
	}
}
