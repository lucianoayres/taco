# 🌮 Taco

![Taco Banner](https://github.com/lucianoayres/taco/blob/main/images/banner_taco.png?raw=true)

## Total Assembly and Concatenation Operator

Welcome to **Taco**, the tool that rolls up all your text files into one deliciously satisfying file. Think of it as a burrito, but for your files! Whether you need to combine your notes, code snippets, or any other text, Taco is here to help with all your concatenation cravings.

## Why Taco?

Because who doesn’t love tacos? But seriously, organizing files can be a mess. With Taco, you can:

-   Concatenate all text files in a directory into one file.
-   Append files to your output without overwriting, because nobody likes to lose their work.

Taco is simple, light, and ready to roll—just like a taco!

## Features

-   🌮 **Concatenate all text files** in the current directory.
-   🚫 **Skips hidden files and directories**, binary files, and the output file.
-   ✨ **Optional Arguments** for directories and output file name.
-   🔄 **Append Mode** to avoid accidental overwrites.

## Project Structure

Here’s the project structure for Taco:

```
/taco
├── go.mod       # Go module file
├── Makefile     # Makefile to simplify commands
├── src          # Directory containing the source code
│   └── main.go  # Main Go file
```

Your main logic for Taco is in the `src/main.go` file, but we will run and build the project from the root directory.

## Getting Started

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

### Default Use (In Current Directory)

Just run **Taco** without any arguments, and it will concatenate all the **text files** in the directory where the command is executed (excluding hidden files, binary files, and itself):

```bash
taco
```

This will create (or append to) a file called `taco.txt`, which will contain:

```
// Filename: example.txt

<Contents of example.txt>

// Filename: notes.md

<Contents of notes.md>
```

### Custom Output File

Want to save your file mashup to something other than the default `taco.txt`? No problem! Just specify your preferred output file name:

```bash
taco -output=my-taco.txt
```

Now everything goes into `my-taco.txt`, all wrapped up like a taco filled with data goodness.

### Custom Directories

You can also tell Taco which directories to collect files from. Taco will process the **text files directly in those directories** (not in subdirectories), and it will exclude hidden and binary files.

```bash
taco /path/to/dir1 /path/to/dir2
```

Or combine it with a custom output:

```bash
taco -output my-taco.txt /path/to/dir1 /path/to/dir2
```

### Skipping Files You Don’t Want

Taco automatically skips:

-   The script itself (because no taco should eat itself).
-   The output file (so it doesn’t end up eating its own leftovers).
-   **Hidden files and directories**.
-   **Binary files** (like images, executables, etc.).
-   **Files in subdirectories** (only processes files in the specified directories).

## Pro Tips

-   **Only Text Files**: Taco includes only text files in the concatenation. It automatically detects text files based on their content, so no need to worry about file extensions.
-   **Hidden Files and Directories**: Hidden files and directories (those starting with a dot `.` like `.git`) are skipped to keep things clean.
-   **Multiple Directories?** Just specify them all in the command, and Taco will grab text files from all specified directories.
-   **Appending?** Run the same command multiple times, and Taco won’t overwrite your carefully crafted file—it’ll just keep adding to it like a buffet plate!

## Makefile Commands

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

## Examples

### Concatenate Text Files in Current Directory

```bash
taco
```

This command concatenates all text files in the current directory into `taco.txt`.

### Concatenate Text Files with Custom Output Filename

```bash
taco -output=my-concatenated-files.txt
```

All text files in the current directory will be concatenated into `my-concatenated-files.txt`.

### Concatenate Text Files from Multiple Directories

```bash
taco -output=my-taco.txt /path/to/dir1 /path/to/dir2
```

This command concatenates all text files from `/path/to/dir1` and `/path/to/dir2` into `my-taco.txt`.

## Limitations

-   **No Subdirectory Traversal**: Taco processes files only in the specified directories and does not traverse subdirectories.
-   **Binary Files Excluded**: Binary files (like images, videos, executables) are automatically excluded to prevent unexpected results.
-   **Hidden Files Skipped**: Files and directories starting with a dot `.` are considered hidden and are skipped.

## Contributions 🍽️

Found a bug? Have a feature idea? Or maybe you just want to share your love of tacos? Feel free to open an issue or create a pull request on [GitHub](https://github.com/lucianoayres/taco). Let’s make Taco even tastier!

### License

This project is licensed under the [MIT License](LICENSE). Eat tacos responsibly.
