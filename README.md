# 🌮 Taco

![Taco Banner](https://github.com/lucianoayres/taco/blob/main/images/banner_taco.png?raw=true)

## Roll up all your text files from directories into one simple text file!

[What's Taco? 🌮](#whats-taco-) · [In Action 🌶️](#in-action-) · [Why Taco? 🤔](#why-taco-) · [Features ✨](#features-) · [Project Structure 📁](#project-structure-) · [Getting Started 🚀](#getting-started-) · [Arguments 📝](#arguments-) · [How to Use Taco 🌮](#how-to-use-taco-) · [Pro Tips 💡](#pro-tips-) · [Makefile Commands 🛠️](#makefile-commands-) · [Examples 📚](#examples-) · [Limitations ⚠️](#limitations-) · [Roadmap 🗺️](#roadmap-) · [Contributions 🍽️](#contributions-) · [License 📄](#license-)

## What's Taco? 🌮

**Taco** is a simple tool that gathers all text files from a directory and its subdirectories, combining them into one file. Ideal for preparing source code or text content for use as prompts with large language models (LLMs), Taco streamlines the process, saving time by automating file gathering.

## In Action 🌶️

![Taco Demo](https://github.com/lucianoayres/taco/blob/main/images/demo_taco.gif?raw=true)

## Why Taco? 🤔

-   🗂️ **Efficient Aggregation**: Automatically collects large numbers of text files.
-   ⏱️ **Saves Time**: Eliminates manual copying, enhancing workflow.
-   🛡️ **Reduces Errors**: Minimizes inconsistencies and mistakes.
-   🤖 **LLM Integration**: Ideal for creating data-rich prompts.
-   📚 **Simplifies Documentation**: Compiles documentation and source code.
-   🚀 **Boosts Productivity**: Streamlines tasks to enhance efficiency.

### Features ✨

-   🌮 **Automatic Gathering**: Collects text files from specified directories and subdirectories.
-   📂 **Recursive Processing**: Seamlessly handles all nested directories.
-   📁 **Flexible Directory and File Selection**: Customizable file and directory filters.
-   🚫 **Skip Hidden and Binary Files**: Keeps output clean.
-   📝 **Detailed Status Updates**: Displays progress and skips details.
-   🔄 **Append Mode**: Adds to existing files without overwriting.
-   ✨ **Customizable Output**: Set custom output names and locations.
-   🎯 **Exclude Files by Pattern**: Exclude files matching specific patterns or regular expressions.

## Project Structure 📁

```
/taco
├── go.mod       # Go module file
├── Makefile     # Makefile to simplify commands
├── src          # Directory containing the source code
│   └── main.go  # Main Go file
```

## Getting Started 🚀

### Installation Options

#### Option 1: Download Pre-built Binary (Linux Only)

1. **Navigate to the [Releases Page](https://github.com/lucianoayres/taco/releases)**.
2. **Download the Latest Release** for Linux.
3. **Make the Binary Executable**:

    ```bash
    chmod +x taco
    ```

4. **Move the Binary to `/usr/local/bin/`**:

    ```bash
    sudo mv taco /usr/local/bin/
    ```

#### Option 2: Build from Source

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/lucianoayres/taco.git
    cd taco
    ```

2. **Initialize the Go Module**:

    ```bash
    go mod tidy
    ```

3. **Build Taco**:

    ```bash
    make build
    ```

4. **Install Taco**:

    ```bash
    sudo make install
    ```

### Running Taco

After setting up Taco, you can run it in one of the following ways:

#### Option 1: Running from the Binary (Preferred)

If you downloaded the binary or built it yourself, you can run Taco directly by simply typing:

```bash
taco
```

This will execute the binary from anywhere if it’s located in your system’s PATH (e.g., `/usr/local/bin/`). Alternatively, if the binary is in your current directory, you can run:

```bash
./taco
```

#### Option 2: Running from the Source Code

To run Taco directly from the source without building, navigate to the project’s root directory and use:

```bash
go run ./src
```

This command will execute the Taco code in the `src` directory, bypassing the need to build the binary first.

## Arguments 📝

Use these optional flags to customize how Taco processes your files:

-   **`-output`**

    -   Specifies the output file name (default: `taco.txt`).

-   **`-include-ext`**

    -   Specify file extensions to include (e.g., `.go,.md`).

-   **`-exclude-ext`**

    -   Specify file extensions to exclude (e.g., `.test,.spec.js`).

-   **`-exclude-dir`**

    -   Exclude root-level directories from processing.

-   **`-include-dir`**

    -   Include specific directories for processing.

-   **`-exclude-file-pattern`**

    -   Exclude files matching specific patterns or regular expressions (e.g., `.*_test\.go$,^LICENSE$`).

-   **`-verbose`**

    -   Enables verbose output for detailed status messages.

---

### Notes

-   **Precedence**: When both `-include-ext` and `-exclude-ext` are used, Taco first filters files by `-include-ext`, then excludes any files that match `-exclude-ext`. The same applies to `-include-dir` and `-exclude-dir`.
-   **Extension Format**: Taco automatically prepends a dot if missing.
-   **Pattern Matching**: The `-exclude-file-pattern` flag uses regular expressions for pattern matching. Ensure patterns are valid and properly escaped.

## How to Use Taco 🌮

### Basic Usage

Run Taco without arguments to concatenate all text files in the current directory and its subdirectories (excluding hidden and binary files):

```bash
taco
```

### Customizing File Selection

#### Specifying Directories to Include or Exclude

-   **Include specific directories**: Use `-include-dir` to process only specified directories. If omitted, Taco processes all subdirectories in the current directory.

    ```bash
    taco -include-dir=src,docs
    ```

-   **Exclude specific directories**: Use `-exclude-dir` to skip certain directories. For instance, to skip `vendor` and `tests`:

    ```bash
    taco -exclude-dir=vendor,tests
    ```

> **Note:** The `-exclude-dir` flag applies only to root-level directories of the specified path.

#### Filtering Files by Extension

Choose files by extension using these flags:

-   **Include specific file types**: Use `-include-ext` to focus on particular file types, such as `.go` and `.md`.

    ```bash
    taco -include-ext=.go,.md
    ```

-   **Exclude certain file types**: Use `-exclude-ext` to omit specific extensions, like `.test` and `.spec.js`.

    ```bash
    taco -exclude-ext=.test,.spec.js
    ```

> **Tip:** Taco first filters files by `-include-ext`, then removes those matching `-exclude-ext`, if both are specified.

#### Excluding Files by Pattern

Use the `-exclude-file-pattern` flag to exclude files matching specific patterns or regular expressions.

-   **Exclude test files and backups**:

    ```bash
    taco -exclude-file-pattern=".*_test\.go$,.*\.bak$"
    ```

-   **Exclude specific files**:

    ```bash
    taco -exclude-file-pattern="^Makefile$,^LICENSE$"
    ```

> **Note:** Patterns are regular expressions. Ensure they are properly quoted and escaped.

### Combining Options

Combine flags to refine file selection. For example:

```bash
taco -output=my-taco.txt -include-dir=src,docs -exclude-dir=vendor,tests -include-ext=.go,.md -exclude-ext=.log,.tmp -exclude-file-pattern=".*_test\.go$,^Makefile$,^LICENSE$"
```

This command combines `.go` and `.md` files from `src` and `docs`, excludes `vendor` and `tests`, omits `.log` and `.tmp` files, and excludes files matching the specified patterns.

## Pro Tips 💡

-   **Recursive Processing**: Taco traverses all subdirectories.
-   **Only Text Files**: Includes files based on content.
-   **Custom Output**: Set output file with `-output`.
-   **Include/Exclude Specific Extensions**: Use `-include-ext` or `-exclude-ext`.
-   **Exclude Directories**: Use `-exclude-dir` to omit directories.
-   **Exclude Files by Pattern**: Use `-exclude-file-pattern` for fine-grained file exclusion.
-   **Append Mode**: Appends new content to the existing output file if it already exists.
-   **Detailed Status**: Verbose mode for skip reasons.

## Makefile Commands 🛠️

Simplify Taco usage with these Makefile commands:

-   **Run Taco**: `make run`
-   **Build Taco**: `make build`
-   **Install Taco**: `sudo make install`
-   **Clean Build Files**: `make clean`
-   **Uninstall Taco**: `sudo make uninstall`
-   **Run Tests**: `make test`
-   **Run Tests with Coverage**: `make test-coverage`
-   **Generate HTML Coverage Report**: `make coverage-html`

## Examples 📚

-   **Concatenate Text Files in Current Directory**:

    ```bash
    taco
    ```

-   **Custom Output Filename**:

    ```bash
    taco -output=my-files.txt
    ```

-   **Specify Directories**:

    ```bash
    taco -output=my-taco.txt -include-dir=src,docs
    ```

-   **Include Only Go and Markdown Files**:

    ```bash
    taco -include-ext=.go,.md
    ```

-   **Exclude Specific Extensions**:

    ```bash
    taco -exclude-ext=.test,.spec.js
    ```

-   **Exclude Files Matching Patterns**:

    ```bash
    taco -exclude-file-pattern=".*_test\.go$,^Makefile$,^LICENSE$"
    ```

-   **Exclude Directories and Set Custom Output**:

    ```bash
    taco -output=final.txt -exclude-dir=vendor,tests -include-dir=src,docs -include-ext=.py,.md -exclude-ext=.log,.tmp -exclude-file-pattern=".*_test\.py$"
    ```

## Limitations ⚠️

-   **Binary Files Excluded**: Binary files are automatically excluded.
-   **Hidden Files Skipped**: Files/directories starting with a dot are skipped.
-   **File Extension Detection**: Relies on extensions for inclusion/exclusion.
-   **Pattern Matching**: Exclusion by pattern uses regular expressions, which require valid syntax.
-   **Conflict Handling**: When using both inclusion and exclusion arguments, overlapping criteria may need careful management.

## Roadmap 🗺️

-   [x] **Launch v1.0**
-   [x] **Implement `-include-ext` feature** (Completed)
-   [x] **Implement `-exclude-ext` feature** (Completed)
-   [x] **Implement `-exclude-dir` feature** (Completed)
-   [x] **Implement `-include-dir` feature** (Completed)
-   [x] **Add regex-based filename exclusion (`-exclude-file-pattern`)** (Completed)
-   [ ] **Add regex-based filename inclusion**
-   [ ] **Support for `.gitignore` files**
-   [ ] **Enhanced error handling and logging**
-   [ ] **Cross-platform binary releases**

## Contributions 🍽️

Found a bug or have a feature request? Open an issue or create a pull request on [GitHub](https://github.com/lucianoayres/taco). Let's make Taco even better!

## License 📄

This project is licensed under the [MIT License](LICENSE). Enjoy your tacos responsibly!
