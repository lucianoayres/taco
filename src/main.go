package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const bufferSize = 512 // Number of bytes to read for content detection

// parseArguments handles the command-line arguments and returns the output filename and directories to process.
func parseArguments() (string, []string) {
	// Define command-line flags
	outputFileName := flag.String("output", "taco.txt", "The output file where the content will be concatenated")
	flag.Parse()

	// Collect directories from arguments or default to current directory
	directories := flag.Args()
	if len(directories) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
		directories = []string{cwd}
	}

	return *outputFileName, directories
}

// getExcludedFiles returns a map containing filenames to exclude (the script itself and the output file).
func getExcludedFiles(outputFileName, scriptFileName string) map[string]struct{} {
	excludedFiles := make(map[string]struct{})
	excludedFiles[scriptFileName] = struct{}{} // Exclude the script/executable file
	excludedFiles[outputFileName] = struct{}{} // Exclude the output file
	return excludedFiles
}

// concatenateFiles processes the directories and writes the content of each text file to the output file.
func concatenateFiles(outputFileName string, directories []string, excludedFiles map[string]struct{}) error {
	// Open or create the output file in append mode
	outputFile, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error creating/opening output file: %v", err)
	}
	defer outputFile.Close()

	for _, dir := range directories {
		err := processDirectory(dir, outputFile, excludedFiles)
		if err != nil {
			return fmt.Errorf("error processing directory %s: %v", dir, err)
		}
	}

	return nil
}

// processDirectory reads files in the directory and processes each text file.
func processDirectory(dir string, outputFile *os.File, excludedFiles map[string]struct{}) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip hidden files
		if isHidden(name) {
			continue
		}

		// Skip excluded files
		if _, excluded := excludedFiles[name]; excluded {
			continue
		}

		// Full path to the file
		path := filepath.Join(dir, name)

		// Check if the file is a text file
		if isTextFile(path) {
			// Write file content to the output file
			err := writeFileContent(outputFile, path, name)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", name, err)
			}
		}
	}

	return nil
}

// writeFileContent reads a file and writes its content to the output file in the specified format.
func writeFileContent(outputFile *os.File, filePath, fileName string) error {
	// Open the file for reading
	inputFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	defer inputFile.Close()

	// Write the filename to the output file
	if _, err := fmt.Fprintf(outputFile, "// Filename: %s\n\n", fileName); err != nil {
		return fmt.Errorf("error writing filename to output file: %v", err)
	}

	// Copy the file content to the output file
	if _, err := io.Copy(outputFile, inputFile); err != nil {
		return fmt.Errorf("error copying content from %s: %v", filePath, err)
	}

	// Write two newlines to separate files
	if _, err := outputFile.WriteString("\n\n"); err != nil {
		return fmt.Errorf("error writing separator to output file: %v", err)
	}

	return nil
}

// isHidden checks if a file is hidden (starts with a dot).
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// isTextFile determines if a file is a text file using net/http.DetectContentType.
func isTextFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		// If we can't open it, assume it's not text
		return false
	}
	defer file.Close()

	// Read the first 512 bytes for content detection
	buffer := make([]byte, bufferSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}
	buffer = buffer[:n]

	contentType := http.DetectContentType(buffer)
	return strings.HasPrefix(contentType, "text/") || contentType == "application/json" || contentType == "application/javascript"
}

func main() {
	// Get the executable's base name to exclude it from concatenation
	scriptFilePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		os.Exit(1)
	}
	scriptFileName := filepath.Base(scriptFilePath)

	// Parse command-line arguments
	outputFileName, directories := parseArguments()

	// Get the list of files to exclude (script and output file)
	excludedFiles := getExcludedFiles(outputFileName, scriptFileName)

	// Concatenate files from the directories
	if err := concatenateFiles(outputFileName, directories, excludedFiles); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Files concatenated successfully into %s\n", outputFileName)
}
