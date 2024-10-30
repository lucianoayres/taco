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

var initialWorkingDir string

// parseArguments handles the command-line arguments and returns the output filename and directories to process.
func parseArguments() (string, []string) {
	// Define command-line flags
	outputFileName := flag.String("output", "taco.txt", "The output file where the content will be concatenated")
	flag.Parse()

	// Collect directories from arguments or default to current directory
	directories := flag.Args()
	if len(directories) == 0 {
		directories = []string{initialWorkingDir}
	}

	return *outputFileName, directories
}

// getExcludedPaths returns a map containing full paths to exclude (the script itself and the output file).
func getExcludedPaths(outputFilePath, scriptFilePath string) map[string]struct{} {
	excludedPaths := make(map[string]struct{})
	excludedPaths[scriptFilePath] = struct{}{}
	excludedPaths[outputFilePath] = struct{}{}
	return excludedPaths
}

// concatenateFiles processes the directories and writes the content of each text file to the output file.
func concatenateFiles(outputFilePath string, directories []string, excludedPaths map[string]struct{}) error {
	var anyFilesProcessed bool = false
	var outputFile *os.File

	defer func() {
		if outputFile != nil {
			outputFile.Close()
		}
	}()

	for _, dir := range directories {
		filesProcessed, err := processDirectory(dir, &outputFile, outputFilePath, excludedPaths)
		if err != nil {
			return fmt.Errorf("error processing directory %s: %v", dir, err)
		}
		if !filesProcessed {
			relativeDir, err := filepath.Rel(initialWorkingDir, dir)
			if err != nil || relativeDir == "." {
				relativeDir = dir
			}
			fmt.Printf("No text files found in %s\n", relativeDir)
		} else {
			anyFilesProcessed = true
		}
	}

	if !anyFilesProcessed {
		fmt.Println("No text files found in any of the directories.")
	} else {
		// Compute relative path of the output file for display
		relativeOutputPath, err := filepath.Rel(initialWorkingDir, outputFilePath)
		if err != nil {
			relativeOutputPath = outputFilePath // Fallback to absolute path
		}
		fmt.Printf("Files concatenated successfully into %s\n", relativeOutputPath)
	}

	return nil
}

// processDirectory recursively reads files in the directory and its subdirectories.
// It returns a bool indicating whether any text files were processed.
func processDirectory(dir string, outputFile **os.File, outputFilePath string, excludedPaths map[string]struct{}) (bool, error) {
	filesProcessed := false

	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("error reading directory %s: %v", dir, err)
	}

	// Track whether any text files were found in subdirectories
	subdirFilesProcessed := false

	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(dir, name)

		// Skip hidden files and directories
		if isHidden(name) {
			continue
		}

		// Skip excluded files and directories based on full path
		if _, excluded := excludedPaths[path]; excluded {
			continue
		}

		if entry.IsDir() {
			// Recursively process subdirectories
			subdirProcessed, err := processDirectory(path, outputFile, outputFilePath, excludedPaths)
			if err != nil {
				return false, err
			}
			if !subdirProcessed {
				relativeDir, err := filepath.Rel(initialWorkingDir, path)
				if err != nil || relativeDir == "." {
					relativeDir = path
				}
				fmt.Printf("No text files found in %s\n", relativeDir)
			} else {
				subdirFilesProcessed = true
			}
		} else {
			// Check if the file is a text file
			if isTextFile(path) {
				// Open output file if not already opened
				if *outputFile == nil {
					var err error
					*outputFile, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return false, fmt.Errorf("error creating/opening output file: %v", err)
					}
				}

				// Compute relative path from initial working directory
				relativePath, err := filepath.Rel(initialWorkingDir, path)
				if err != nil {
					relativePath = path // Fallback to full path if cannot compute relative path
				}

				// Processing status in a single line
				fmt.Printf("Processing %s ... ", relativePath)

				// Write file content to the output file
				err = writeFileContent(*outputFile, path, relativePath)
				if err != nil {
					fmt.Printf("Error\n")
					fmt.Printf("Error processing file %s: %v\n", relativePath, err)
				} else {
					// Indicate completion on the same line
					fmt.Printf("Done\n")
				}

				filesProcessed = true
			}
		}
	}

	// If no files were processed in this directory or its subdirectories
	if !filesProcessed && !subdirFilesProcessed {
		return false, nil
	}

	return true, nil
}

// writeFileContent reads a file and writes its content to the output file in the specified format.
func writeFileContent(outputFile *os.File, filePath, relativePath string) error {
	// Open the file for reading
	inputFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	defer inputFile.Close()

	// Write the file path to the output file
	if _, err := fmt.Fprintf(outputFile, "// File: %s\n\n", relativePath); err != nil {
		return fmt.Errorf("error writing file path to output file: %v", err)
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

// isHidden checks if a file or directory is hidden (starts with a dot).
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

	// Handle empty files as text files
	if n == 0 {
		return true
	}

	contentType := http.DetectContentType(buffer)
	return strings.HasPrefix(contentType, "text/") ||
		contentType == "application/json" ||
		contentType == "application/javascript" ||
		contentType == "application/xml"
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var err error
	// Get the initial working directory
	initialWorkingDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting current working directory: %v", err)
	}

	// Get the executable's full path to exclude it from concatenation
	scriptFilePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Error getting executable path: %v", err)
	}

	// Parse command-line arguments
	outputFileName, directories := parseArguments()

	// Get the absolute path of the output file
	outputFilePath, err := filepath.Abs(outputFileName)
	if err != nil {
		return fmt.Errorf("Error getting absolute path of output file: %v", err)
	}

	// Get the list of paths to exclude (script and output file)
	excludedPaths := getExcludedPaths(outputFilePath, scriptFilePath)

	// Concatenate files from the directories
	if err := concatenateFiles(outputFilePath, directories, excludedPaths); err != nil {
		return err
	}
	return nil
}
