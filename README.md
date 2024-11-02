# 🌮 Taco

![Taco Banner](https://github.com/lucianoayres/taco/blob/main/images/banner_taco.png?raw=true)

## Roll up all your text files from directories into one simple text file!

[What's Taco? 🌮](#whats-taco-) · [In Action 🌶️](#in-action-) · [Why Taco? 🤔](#why-taco-) · [Features ✨](#features-) · [Project Structure 📁](#project-structure-) · [Getting Started 🚀](#getting-started-) · [How to Use Taco 🌮](#how-to-use-taco-) · [Pro Tips 💡](#pro-tips-) · [Makefile Commands 🛠️](#makefile-commands-) · [Examples 📚](#examples-) · [Limitations ⚠️](#limitations-) · [Roadmap 🗺️](#roadmap-) · [Contributions 🍽️](#contributions-) · [License 📄](#license-)

## What's Taco? 🌮

**Taco** is a simple tool that gathers all the text files from a directory and its subdirectories and combines them into a single text file. It’s perfect for pulling together source code or text content to use as a prompt with large language models (LLMs), making the process easier and more efficient!

With Taco, you can forget about manually copying and pasting individual files. Just one command collects all your text files, like wrapping them into one big burrito! 😄

## In Action 🌶️

![Taco Demo](https://github.com/lucianoayres/taco/blob/main/images/demo_taco.gif?raw=true)

## Why Taco? 🤔

Working with large collections of text files can be cumbersome, especially when preparing data for large language models or consolidating documentation and source code. Manually aggregating these files is not only time-consuming but also prone to errors and inconsistencies. **Taco** streamlines this process by automating the collection and consolidation of text files, ensuring a clean and efficient workflow. Whether you're a developer looking to create comprehensive prompts for LLMs or need to compile documentation from various sources, Taco simplifies your tasks and boosts your productivity.

## Features ✨

-   🌮 **Automatic Text File Gathering**: Collects all text files from specified directories and their subdirectories.
-   📂 **Recursive Traversal**: Processes all nested directories seamlessly.
-   🚫 **Exclusion of Hidden and Binary Files**: Keeps your output clean by skipping unnecessary files.
-   📁 **Include Specific File Types**: Use the `-include-ext` flag to specify which file extensions to include.
-   ❌ **Exclude Specific File Types**: Use the `-exclude-ext` flag to specify which file extensions to exclude.
-   📝 **Status Messages**: Provides clear progress updates during processing.
-   🔄 **Append Mode**: Prevents accidental overwrites by appending to existing output files.
-   ✨ **Customizable Output**: Specify output file names and directories as needed.

## Project Structure 📁

Here’s the project structure for Taco:

```
/taco
├── go.mod       # Go module file
├── Makefile     # Makefile to simplify commands
├── src          # Directory containing the source code
│   └── main.go  # Main Go file
```

Your main logic for Taco is in the `src/main.go` file, but we will run and build the project from the root directory.

## Getting Started 🚀

### Installation Options

You have two options to install Taco:

#### **Option 1: Download Pre-built Binary (Linux Only)**

If you prefer not to build Taco from source, you can download the pre-built binary for Linux:

1. **Navigate to the Releases Page:**

    Visit the [Taco Releases](https://github.com/lucianoayres/taco/releases) page on GitHub.

2. **Download the Latest Release:**

    Find the latest release and download the Linux binary (`taco`).

3. **Make the Binary Executable:**

    After downloading, make the binary executable:

    ```bash
    chmod +x taco
    ```

4. **Move the Binary to Your Local Bin Directory:**

    To use Taco from anywhere, move it to `/usr/local/bin/`:

    ```bash
    sudo mv taco /usr/local/bin/
    ```

    Now, you can run Taco by simply typing:

    ```bash
    taco
    ```

    **Note:** The pre-built binary is **Linux-only**. Ensure you're using a Linux environment to utilize the downloaded binary.

#### **Option 2: Build from Source**

If you prefer to build Taco from source, follow these steps:

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/lucianoayres/taco.git
    cd taco
    ```

2. **Initialize the Go Module:**

    Ensure you have Go installed, then run:

    ```bash
    go mod tidy
    ```

    This will download any required dependencies.

3. **Build Taco:**

    You can build Taco into an executable using the Makefile or directly with `go build`.

    - **Using Makefile:**

        ```bash
        make build
        ```

        This will create an executable file called `taco` in the root directory.

    - **Using `go build`:**

        ```bash
        go build -o taco ./src
        ```

        This also creates the `taco` executable in the root directory.

4. **Install Taco:**

    To install Taco so you can use it from anywhere, move the executable to your local bin directory:

    ```bash
    sudo make install
    ```

    This will build the executable and move it to `/usr/local/bin/`.

    Alternatively, if you built it manually:

    ```bash
    sudo mv taco /usr/local/bin/
    ```

    Now, you can run the program by executing:

    ```bash
    taco
    ```

### Running Taco

To run Taco, you can execute the script directly from the root directory of the project.

#### Option 1: Run Taco from the Root Directory

Simply run:

```bash
go run ./src
```

#### Option 2: Use the Makefile

Use the provided Makefile to simplify commands:

```bash
make run
```

### Building Taco

You can build Taco into an executable. To build the project from the root directory, use:

```bash
go build -o taco ./src
```

Alternatively, use the Makefile:

```bash
make build
```

This will create an executable file called `taco` in the root directory. You can then run the program by executing:

```bash
./taco
```

### Installing Taco

To install Taco so you can use it from anywhere, move the executable to your local bin directory:

```bash
sudo make install
```

This will build the executable and move it to `/usr/local/bin/`.

## How to Use Taco 🌮

### Default Use (In Current Directory and Subdirectories)

Run **Taco** without any arguments, and it will recursively concatenate all the **text files** in the current directory and its subdirectories (excluding hidden and binary files, and itself):

```bash
taco
```

This will create (or append to) a file called `taco.txt`, which will contain:

```
# File: README.md

<Contents of README.md>

# File: src/main.go

<Contents of main.go>

# File: docs/setup/guide.txt

<Contents of guide.txt>
```

### Custom Output File

Specify your preferred output file name:

```bash
taco -output=my-taco.txt
```

Now everything goes into `my-taco.txt`.

### Custom Directories

Specify directories to collect files from. Taco will recursively process the **text files in those directories and all their subdirectories**:

```bash
taco /path/to/dir1 /path/to/dir2
```

Or combine it with a custom output:

```bash
taco -output=my-taco.txt /path/to/dir1 /path/to/dir2
```

### Including Specific File Extensions

Use the `-include-ext` flag to specify which file extensions to include in the concatenation process. This allows you to focus on specific types of files, such as only Go source files and Markdown documentation.

**Usage Example:**

```bash
taco -include-ext=.go,.md
```

This command will concatenate only `.go` and `.md` files from the specified directories.

**Additional Examples:**

-   **Including Extensions Without Leading Dots:**

    ```bash
    taco -include-ext=go,md
    ```

    Taco will automatically prepend a dot to extensions without a leading dot.

-   **Including Extensions with Mixed Case and Spaces:**

    ```bash
    taco -include-ext=".Go, .MD , .Txt"
    ```

    Taco handles case insensitivity and trims whitespace around extensions.

-   **Combining `-include-ext` with Other Flags:**

    ```bash
    taco -output=combined.txt -include-ext=.js,.json /path/to/dir1 /path/to/dir2
    ```

### Excluding Specific File Extensions

Use the new `-exclude-ext` flag to specify which file extensions to exclude from the concatenation process. This enables you to omit certain file types that are not needed in the final output.

**Usage Example:**

```bash
taco -exclude-ext=.test,.spec.js
```

This command will exclude all `.test` and `.spec.js` files from being concatenated.

**Additional Examples:**

-   **Excluding Extensions Without Leading Dots:**

    ```bash
    taco -exclude-ext=test,spec.js
    ```

    Taco will automatically prepend a dot to extensions without a leading dot.

-   **Excluding Extensions with Mixed Case and Spaces:**

    ```bash
    taco -exclude-ext=".Test, .SPEC.JS , .Tmp"
    ```

    Taco handles case insensitivity and trims whitespace around extensions.

-   **Combining `-exclude-ext` with Other Flags:**

    ```bash
    taco -output=filtered.txt -exclude-ext=.log,.tmp /path/to/dir1 /path/to/dir2
    ```

### Including and Excluding Specific File Extensions

You can combine both `-include-ext` and `-exclude-ext` flags to fine-tune which files are included or excluded in the concatenation process.

**Usage Example:**

```bash
taco -include-ext=.go,.md -exclude-ext=.test,.spec.js
```

This command will include only `.go` and `.md` files while excluding any `.test` and `.spec.js` files from those included types.

**Important Notes:**

-   **Precedence:** When both `-include-ext` and `-exclude-ext` are used, Taco first filters files based on the `-include-ext` list and then excludes any files that match the `-exclude-ext` list.
-   **Case Insensitivity:** Both flags handle case insensitivity, so `.GO` is treated the same as `.go`.
-   **Whitespace Handling:** Whitespace around extensions is automatically trimmed.

### Skipping Files You Don’t Want

Taco automatically skips:

-   The script itself.
-   The output file.
-   **Hidden files and directories** (those starting with a dot `.`).
-   **Binary files** (images, executables, etc.).

### Status Messages

As Taco processes your files, it provides status messages:

```
Processing LICENSE ... Done
Processing README.md ... Done
Processing src/main.go ... Done
No text files found in docs/empty_folder
Files concatenated successfully into taco.txt
```

-   **Processing [file] ... Done**: File processed successfully.
-   **No text files found in [directory]**: Directory contains no text files.

## Pro Tips 💡

-   **Recursive Processing**: Taco automatically traverses all subdirectories.
-   **Only Text Files**: Includes only text files based on content, not file extension.
-   **Including Specific Extensions**: Use the `-include-ext` flag to include only files with certain extensions.
-   **Excluding Specific Extensions**: Use the `-exclude-ext` flag to omit files with certain extensions.
-   **Hidden Files and Directories**: Skipped if starting with a dot `.`.
-   **Multiple Directories**: Specify multiple directories to process files from all of them.
-   **Appending**: Taco appends to existing files unless you delete the output file first.
-   **Check for Empty Directories**: Taco informs you about directories without text files.

## Makefile Commands 🛠️

Simplify using Taco with the **Makefile** commands:

-   **Run Taco**:

    ```bash
    make run
    ```

-   **Build Taco**:

    ```bash
    make build
    ```

-   **Install Taco**:

    ```bash
    sudo make install
    ```

-   **Clean Build Files**:

    ```bash
    make clean
    ```

-   **Uninstall Taco**:

    ```bash
    sudo make uninstall
    ```

-   **Run Tests**:

    ```bash
    make test
    ```

-   **Run Tests with Coverage**:

    ```bash
    make test-coverage
    ```

-   **Generate HTML Coverage Report**:

    ```bash
    make coverage-html
    ```

## Examples 📚

### Concatenate Text Files in Current Directory and Subdirectories

```bash
taco
```

### Concatenate Text Files with Custom Output Filename

```bash
taco -output=my-concatenated-files.txt
```

### Concatenate Text Files from Multiple Directories

```bash
taco -output=my-taco.txt /path/to/dir1 /path/to/dir2
```

### Concatenate Only Go and Markdown Files

```bash
taco -include-ext=.go,.md
```

### Concatenate Only Go and Markdown Files with Custom Output

```bash
taco -output=code_docs.txt -include-ext=.go,.md /path/to/code /path/to/docs
```

### Excluding Specific File Extensions

```bash
taco -exclude-ext=.test,.spec.js
```

### Excluding Extensions Without Leading Dots

```bash
taco -exclude-ext=test,spec.js
```

### Excluding Extensions with Mixed Case and Spaces

```bash
taco -exclude-ext=".Test, .SPEC.JS , .Tmp"
```

### Combining `-include-ext` and `-exclude-ext` Flags

```bash
taco -include-ext=.go,.md -exclude-ext=.test,.spec.js
```

## Limitations ⚠️

-   **Binary Files Excluded**: Binary files are automatically excluded.
-   **Hidden Files Skipped**: Files and directories starting with a dot `.` are skipped.
-   **No Exclusion Flags for Directories**: Currently, there are no flags to exclude specific directories.
-   **File Extension Detection**: While `-include-ext` and `-exclude-ext` allow focusing on specific file types, they rely on file extensions to determine inclusion/exclusion, which may not always reflect the file's actual content type.
-   **Conflict Handling**: When both `-include-ext` and `-exclude-ext` are used, files matching both criteria may lead to unexpected behaviors if not carefully managed.

## Roadmap 🗺️

-   [x] **Launch v1.0**
-   [x] **Implement `-include-ext` feature** (Completed)
-   [x] **Implement `-exclude-ext` feature** (Completed)
-   [ ] **Add support for `.gitignore` files**
-   [ ] **Add argument to exclude directories**
-   [ ] **Add argument to exclude files using regex**
-   [ ] **Enhance error handling and logging**
-   [ ] **Support for additional file content types**
-   [ ] **Cross-platform binary releases**

## Contributions 🍽️

Found a bug or have a feature request? Open an issue or create a pull request on [GitHub](https://github.com/lucianoayres/taco). Let's make Taco even better!

### License 📄

This project is licensed under the [MIT License](LICENSE). Enjoy your tacos responsibly!
