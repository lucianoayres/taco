// File: src/main.go

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const bufferSize = 512 // Number of bytes to read for content detection

var initialWorkingDir string

// parseArguments handles the command-line arguments and returns the output filename, directories to process,
// included extensions, excluded extensions, excluded patterns, excluded directories, verbosity flag, and an error if any.
func parseArguments() (string, []string, []string, []string, []string, []string, bool, error) {
	// Define command-line flags
	outputFileName := flag.String("output", "taco.txt", "The output file where the content will be concatenated")
	includeExt := flag.String("include-ext", "", "Comma-separated list of file extensions to include (e.g., .go,.md)")
	excludeExt := flag.String("exclude-ext", "", "Comma-separated list of file extensions to exclude (e.g., .test,.spec.js)")
	excludeDir := flag.String("exclude-dir", "", "Comma-separated list of directories to exclude (e.g., vendor,tests)")
	includeDir := flag.String("include-dir", "", "Comma-separated list of directories to include (e.g., src,docs,images). If not provided, the current directory and all its subdirectories will be processed.")
	excludeFilePattern := flag.String("exclude-file-pattern", "", "Comma-separated list of file patterns or regular expressions to exclude files")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	var directories []string

	if *includeDir != "" {
		// Parse the include-dir flag into directories
		splitDir := strings.Split(*includeDir, ",")
		for _, dir := range splitDir {
			trimmedDir := strings.TrimSpace(dir)
			if trimmedDir != "" {
				directories = append(directories, trimmedDir)
			}
		}
	} else {
		// If include-dir is not provided, process the current directory and all its subdirectories
		directories = append(directories, ".")
	}

	// Parse the include extensions
	var includeExtensions []string
	if *includeExt != "" {
		splitExt := strings.Split(*includeExt, ",")
		for _, ext := range splitExt {
			trimmedExt := strings.TrimSpace(ext)
			if trimmedExt != "" {
				// Ensure extensions start with a dot
				if !strings.HasPrefix(trimmedExt, ".") {
					trimmedExt = "." + trimmedExt
				}
				includeExtensions = append(includeExtensions, strings.ToLower(trimmedExt))
			}
		}
	}

	// Parse the exclude extensions
	var excludeExtensions []string
	if *excludeExt != "" {
		splitExt := strings.Split(*excludeExt, ",")
		for _, ext := range splitExt {
			trimmedExt := strings.TrimSpace(ext)
			if trimmedExt != "" {
				// Ensure extensions start with a dot
				if !strings.HasPrefix(trimmedExt, ".") {
					trimmedExt = "." + trimmedExt
				}
				excludeExtensions = append(excludeExtensions, strings.ToLower(trimmedExt))
			}
		}
	}

	// Parse the exclude file patterns
	var excludePatterns []string
	if *excludeFilePattern != "" {
		splitPatterns := strings.Split(*excludeFilePattern, ",")
		for _, pattern := range splitPatterns {
			trimmedPattern := strings.TrimSpace(pattern)
			if trimmedPattern != "" {
				excludePatterns = append(excludePatterns, trimmedPattern)
			}
		}
	}

	// Parse the exclude directories
	var excludedDirectories []string
	if *excludeDir != "" {
		splitDir := strings.Split(*excludeDir, ",")
		for _, dir := range splitDir {
			trimmedDir := strings.TrimSpace(dir)
			if trimmedDir != "" {
				excludedDirectories = append(excludedDirectories, trimmedDir)
			}
		}
	}

	return *outputFileName, directories, includeExtensions, excludeExtensions, excludePatterns, excludedDirectories, *verbose, nil
}

// getExcludedPaths returns a map containing full paths to exclude (the script itself and the output file).
func getExcludedPaths(outputFilePath, scriptFilePath string) map[string]struct{} {
	excludedPaths := make(map[string]struct{})
	excludedPaths[scriptFilePath] = struct{}{}
	excludedPaths[outputFilePath] = struct{}{}
	return excludedPaths
}

// concatenateFiles processes the directories and writes the content of each text file to the output file.
func concatenateFiles(outputFilePath string, directories []string, excludedPaths map[string]struct{}, excludedDirs map[string]struct{}, includeExts, excludeExts []string, excludePatterns []*regexp.Regexp, verbose bool) error {
	var anyFilesProcessed bool = false
	var outputFile *os.File

	defer func() {
		if outputFile != nil {
			outputFile.Close()
		}
	}()

	for _, dir := range directories {
		// Resolve the absolute path of the directory
		absDir, err := filepath.Abs(filepath.Join(initialWorkingDir, dir))
		if err != nil {
			return fmt.Errorf("error resolving absolute path of directory %s: %v", dir, err)
		}

		// Check if the directory exists
		info, err := os.Stat(absDir)
		if err != nil {
			if os.IsNotExist(err) {
				if verbose {
					fmt.Printf("Directory does not exist: %s\n", absDir)
				}
				continue
			}
			return fmt.Errorf("error accessing directory %s: %v", absDir, err)
		}
		if !info.IsDir() {
			if verbose {
				fmt.Printf("Not a directory, skipping: %s\n", absDir)
			}
			continue
		}

		filesProcessed, err := processDirectory(absDir, &outputFile, outputFilePath, excludedPaths, excludedDirs, includeExts, excludeExts, excludePatterns, verbose)
		if err != nil {
			return fmt.Errorf("error processing directory %s: %v", dir, err)
		}
		if !filesProcessed && verbose {
			relativeDir, err := filepath.Rel(initialWorkingDir, absDir)
			if err != nil || relativeDir == "." {
				relativeDir = dir
			}
			fmt.Printf("No text files found in %s\n", relativeDir)
		} else if filesProcessed {
			anyFilesProcessed = true
		}
	}

	if !anyFilesProcessed && verbose {
		fmt.Println("No text files found in any of the directories.")
	} else if anyFilesProcessed {
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
func processDirectory(dir string, outputFile **os.File, outputFilePath string, excludedPaths map[string]struct{}, excludedDirs map[string]struct{}, includeExts, excludeExts []string, excludePatterns []*regexp.Regexp, verbose bool) (bool, error) {
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
			if verbose {
				fmt.Printf("Skipping excluded path: %s\n", path)
			}
			continue
		}

		// Check if the current directory is in the excluded directories
		if entry.IsDir() {
			relPath, err := filepath.Rel(initialWorkingDir, path)
			if err != nil {
				relPath = path // Fallback to absolute path
			}
			// Normalize the relative path for comparison
			normalizedRelPath := filepath.ToSlash(relPath)
			if _, excluded := excludedDirs[normalizedRelPath]; excluded {
				if verbose {
					fmt.Printf("Skipping excluded directory: %s\n", relPath)
				}
				continue
			}

			// Recursively process subdirectories
			subdirProcessed, err := processDirectory(path, outputFile, outputFilePath, excludedPaths, excludedDirs, includeExts, excludeExts, excludePatterns, verbose)
			if err != nil {
				return false, err
			}
			if !subdirProcessed && verbose {
				relativeDir, err := filepath.Rel(initialWorkingDir, path)
				if err != nil || relativeDir == "." {
					relativeDir = path
				}
				fmt.Printf("No text files found in %s\n", relativeDir)
			} else if subdirProcessed {
				subdirFilesProcessed = true
			}
		} else {
			// Determine relative path once for both processing and exclusion messages
			relativePath, err := filepath.Rel(initialWorkingDir, path)
			if err != nil {
				relativePath = path // Fallback to absolute path if cannot compute relative path
			}

			// Check if the file matches any of the exclude patterns
			if matchesExcludePatterns(name, excludePatterns) {
				if verbose {
					fmt.Printf("Skipping file %s: matches exclude pattern\n", relativePath)
				}
				continue
			}

			// Check if the file should be included based on extensions
			if shouldIncludeFile(path, includeExts, excludeExts) && isTextFile(path) {
				// Open output file if not already opened
				if *outputFile == nil {
					var err error
					*outputFile, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return false, fmt.Errorf("error creating/opening output file: %v", err)
					}
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
			} else {
				// Determine reason for exclusion
				ext := strings.ToLower(filepath.Ext(path))
				skipped := false

				if len(includeExts) > 0 {
					// If includeExts is specified and file is not included
					included := false
					for _, includeExt := range includeExts {
						if ext == includeExt {
							included = true
							break
						}
					}
					if !included {
						if verbose {
							fmt.Printf("Skipping file %s: does not match include extensions\n", relativePath)
						}
						skipped = true
					}
				}
				if !skipped && len(excludeExts) > 0 {
					// If excludeExts is specified and file is excluded
					excluded := false
					for _, excludeExt := range excludeExts {
						if ext == excludeExt {
							excluded = true
							break
						}
					}
					if excluded {
						if verbose {
							fmt.Printf("Skipping file %s: excluded by extension %s\n", relativePath, ext)
						}
						skipped = true
					}
				}
				// If neither include nor exclude reason matches, no message is printed
			}
		}
	}

	// If no files were processed in this directory or its subdirectories
	if !filesProcessed && !subdirFilesProcessed {
		return false, nil
	}

	return true, nil
}

// matchesExcludePatterns checks if the filename matches any of the exclude patterns.
func matchesExcludePatterns(filename string, excludePatterns []*regexp.Regexp) bool {
	for _, re := range excludePatterns {
		if re.MatchString(filename) {
			return true
		}
	}
	return false
}

// shouldIncludeFile determines if a file should be included based on the provided include and exclude extensions.
// If includeExts is empty, it includes all text files except those in excludeExts.
func shouldIncludeFile(filePath string, includeExts, excludeExts []string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check include extensions
	if len(includeExts) > 0 {
		included := false
		for _, includeExt := range includeExts {
			if ext == includeExt {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	// Check exclude extensions
	if len(excludeExts) > 0 {
		for _, excludeExt := range excludeExts {
			if ext == excludeExt {
				return false
			}
		}
	}

	return true
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

	// Write one newline to separate files
	if _, err := outputFile.WriteString("\n"); err != nil {
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
	outputFileName, directories, includeExts, excludeExts, excludePatterns, excludeDirs, verbose, err := parseArguments()
	if err != nil {
		return err
	}

	// Get the absolute path of the output file
	outputFilePath, err := filepath.Abs(outputFileName)
	if err != nil {
		return fmt.Errorf("Error getting absolute path of output file: %v", err)
	}

	// Get the list of paths to exclude (script and output file)
	excludedPaths := getExcludedPaths(outputFilePath, scriptFilePath)

	// Prepare excluded directories map
	excludedDirsMap := make(map[string]struct{})
	for _, dir := range excludeDirs {
		// Normalize directory paths to use forward slashes
		normalizedDir := filepath.ToSlash(dir)
		excludedDirsMap[normalizedDir] = struct{}{}
	}

	// Compile exclude patterns into regular expressions
	var excludeRegexps []*regexp.Regexp
	for _, pattern := range excludePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("Invalid exclude-file-pattern %q: %v", pattern, err)
		}
		excludeRegexps = append(excludeRegexps, re)
	}

	// Concatenate files from the directories
	if err := concatenateFiles(outputFilePath, directories, excludedPaths, excludedDirsMap, includeExts, excludeExts, excludeRegexps, verbose); err != nil {
		return err
	}
	return nil
}
