# 🌮 Taco

![Taco Banner](https://github.com/lucianoayres/taco/blob/main/images/banner_taco.png?raw=true)

## Roll up all your text files from directories into one simple text file!

[What's Taco? 🌮](#whats-taco-) · [Why Taco? 🤔](#why-taco-) · [Features ✨](#features-) · [Project Structure 📁](#project-structure-) · [Getting Started 🚀](#getting-started-) · [How to Use Taco 🌮](#how-to-use-taco-) · [Pro Tips 💡](#pro-tips-) · [Makefile Commands 🛠️](#makefile-commands-) · [Examples 📚](#examples-) · [Limitations ⚠️](#limitations-) · [Contributions 🍽️](#contributions-) · [License 📄](#license-)

## What's Taco? 🌮

Taco is a tool designed to make your workflow with large language models (LLMs) easier. Instead of manually copying and pasting individual project files into a prompt, Taco rolls up all your text files into one neat file, so you can focus on what matters. Think of it like a burrito for your files—whether it’s notes, code snippets, or any other text, Taco simplifies the process by creating a single file with everything you need.

## Why Taco? 🤔

Because manually handling files for LLM prompts can be time-consuming and messy. With Taco, you can:

-   🌮 **Concatenate all text files** in a directory and its subdirectories into one comprehensive file for easy prompting.
-   📂 **Recursively process subdirectories**, gathering all your important text files automatically.
-   🚀 **Display processing status** for each file, giving you clear progress updates.
-   ✨ **Optional arguments** allow you to specify which directories and file names to use.
-   🔄 **Append mode** helps you avoid accidental overwrites, so you can keep adding to your project file seamlessly.

Taco is your go-to solution for organizing and prepping files, all with the efficiency of a single command!

## Features ✨

-   🌮 **Concatenate all text files** in the specified directories and their subdirectories to create a unified prompt file.
-   📂 **Recursive traversal** ensures all nested text files are included.
-   🚫 **Exclusion of binary and hidden files**, keeping your output lean and clean.
-   🔄 **Append mode** prevents overwriting, only growing the output file when needed.
-   📝 **Status messages** that show file-by-file progress.
-   📁 **Directory reporting** for directories without text files, ensuring no data is missed.

Taco is your go-to solution for organizing and prepping files, all with the efficiency of a single command!

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

### 1. Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/lucianoayres/taco.git
cd taco
```

### 2. Initialize the Go Module

If you're cloning the project for the first time, ensure you have Go installed, then run:

```bash
go mod tidy
```

This will download any required dependencies.

### 3. Running Taco

To run Taco, you can execute the script directly from the root directory of the project.

#### Option 1: Run Taco from the Root Directory

Simply run:

```bash
go run ./src
```

#### Option 2: Use the Makefile

We've provided a Makefile to simplify commands. You can run:

```bash
make run
```

### 4. Building Taco

You can also build Taco into an executable. To build the project from the root directory, use this command:

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

### 5. Installing Taco

To install Taco so you can use it from anywhere, you can move the executable to your local bin directory:

```bash
sudo make install
```

This will build the executable and move it to `/usr/local/bin/`.

## How to Use Taco 🌮

### Default Use (In Current Directory and Subdirectories)

Just run **Taco** without any arguments, and it will recursively concatenate all the **text files** in the directory where the command is executed (excluding hidden files, binary files, and itself):

```bash
taco
```

This will create (or append to) a file called `taco.txt`, which will contain:

```
// File: README.md

<Contents of README.md>

// File: src/main.go

<Contents of main.go>

// File: docs/setup/guide.txt

<Contents of guide.txt>
```

### Custom Output File

Want to save your file mashup to something other than the default `taco.txt`? No problem! Just specify your preferred output file name:

```bash
taco -output=my-taco.txt
```

Now everything goes into `my-taco.txt`, all wrapped up like a taco filled with data goodness.

### Custom Directories

You can also tell Taco which directories to collect files from. Taco will recursively process the **text files in those directories and all their subdirectories**, excluding hidden and binary files.

```bash
taco /path/to/dir1 /path/to/dir2
```

Or combine it with a custom output:

```bash
taco -output=my-taco.txt /path/to/dir1 /path/to/dir2
```

### Skipping Files You Don’t Want

Taco automatically skips:

-   The script itself (because no taco should eat itself).
-   The output file (so it doesn’t end up eating its own leftovers).
-   **Hidden files and directories** (those starting with a dot `.`).
-   **Binary files** (like images, executables, etc.).

### Status Messages

As Taco processes your files, it provides informative status messages:

```bash
Processing LICENSE ... Done
Processing README.md ... Done
Processing src/main.go ... Done
No text files found in docs/empty_folder
Files concatenated successfully into taco.txt
```

-   **Processing [file] ... Done**: Indicates that a file has been processed successfully.
-   **No text files found in [directory]**: Informs you when a directory (or subdirectory) contains no text files.

## Pro Tips 💡

-   **Recursive Processing**: Taco automatically traverses all subdirectories to find text files.
-   **Only Text Files**: Taco includes only text files in the concatenation. It automatically detects text files based on their content, so no need to worry about file extensions.
-   **Hidden Files and Directories**: Files and directories starting with a dot `.` are considered hidden and are skipped.
-   **Multiple Directories?** Just specify them all in the command, and Taco will grab text files from all specified directories and their subdirectories.
-   **Appending?** Run the same command multiple times, and Taco won’t overwrite your carefully crafted file—it’ll just keep adding to it like a buffet plate!
-   **Check for Empty Directories**: Taco informs you about directories without text files, so you can keep your folders tidy.

## Makefile Commands 🛠️

To simplify using Taco, we've provided a **Makefile** with handy commands:

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

## Examples 📚

### Concatenate Text Files in Current Directory and Subdirectories

```bash
taco
```

This command concatenates all text files in the current directory and its subdirectories into `taco.txt`.

### Concatenate Text Files with Custom Output Filename

```bash
taco -output=my-concatenated-files.txt
```

All text files in the current directory and its subdirectories will be concatenated into `my-concatenated-files.txt`.

### Concatenate Text Files from Multiple Directories

```bash
taco -output=my-taco.txt /path/to/dir1 /path/to/dir2
```

This command recursively concatenates all text files from `/path/to/dir1` and `/path/to/dir2` into `my-taco.txt`.

### Exclude Specific Files or Directories

While Taco automatically skips hidden and binary files, if you want to exclude specific files or directories, you can reorganize your folders or temporarily rename files. Future versions may include flags for exclusion.

## Limitations ⚠️

-   **Binary Files Excluded**: Binary files (like images, videos, executables) are automatically excluded to prevent unexpected results.
-   **Hidden Files Skipped**: Files and directories starting with a dot `.` are considered hidden and are skipped.
-   **No Exclusion Flags**: Currently, there are no flags to exclude specific files or directories from processing.

## Contributions 🍽️

Found a bug? Have a feature idea? Or maybe you just want to share your love of tacos? Feel free to open an issue or create a pull request on [GitHub](https://github.com/lucianoayres/taco). Let’s make Taco even tastier!

### License 📄

This project is licensed under the [MIT License](LICENSE). Eat tacos responsibly.
