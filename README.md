# DML - Display Markdown & LaTeX

DML is a command-line tool that reads text from standard input, processes basic Markdown (bold/italic) and LaTeX math expressions (both inline and display), and prints the result to standard output. LaTeX math is rendered as images using the Kitty terminal graphics protocol.

## Code Structure

The codebase is organized in a modular structure:

- `cmd/dml/` - Main application entry point and CLI processing
- `internal/` - Core implementation packages:
  - `color/` - Color processing and management
  - `latex/` - LaTeX rendering and document generation
  - `markdown/` - Markdown processing and formatting
  - `regex/` - Regular expression patterns for text analysis
  - `terminal/` - Terminal output and Kitty protocol implementation

## Features

*   Renders inline LaTeX math (`$formula$`) as images.
*   Renders display LaTeX math (`$$formula$$`) as images.
*   Converts Markdown `*italic*` / `_italic_` to italicized text (requires terminal/font support for ANSI italics).
*   Converts Markdown `**bold**` / `__bold__` to bold text.
*   Customizable text color for rendered LaTeX images.
*   Passes through unrecognized Markdown and other text.

## Prerequisites

Before using DML, you need the following installed on your system:

1.  **Go**: Version 1.18 or higher (to build the tool).
2.  **pdflatex**: For compiling LaTeX expressions. This is part of standard TeX distributions like TeX Live or MiKTeX.
3.  **convert**: For converting the PDF output from LaTeX into PNG images. This is part of the ImageMagick suite.
4.  **A Kitty-compatible terminal**: Required to display the inline images rendered by DML.

## Installation

1.  **Clone the Repository (if applicable)**:
    If you have this project in a Git repository:
    ```bash
    git clone https://github.com/JamieLittle16/dml
    cd dml
    ```
    If you just have the source files, navigate to the `Projects/dml` directory.

2.  **Build the Binary**:
    From the `Projects/dml` directory:
    ```bash
    ./build.sh
    ```
    This will create an executable file named `dml` in the current directory.

3.  **Install the Binary and Man Page (System-Wide)**:
    To make `dml` available globally and install its man page, run the following commands from the `Projects/dml` directory:

    ```bash
    ./build.sh install
    ```
    This will install the binary to `/usr/local/bin/dml` and the man page to `/usr/local/share/man/man1/dml.1.gz`.

    *Note: This requires sudo access. After installation, you might need to run `sudo mandb` (or `rehash` in some shells) for the system to recognize the new man page and command.*

4.  **User-Specific Installation (Optional)**:
    Alternatively, you can install it to a user-specific directory like `$HOME/.local/bin` (ensure this is in your PATH):
    ```bash
    ./build.sh install local
    ```
    This will install the binary to `$HOME/.local/bin/dml` and the man page to `$HOME/.local/share/man/man1/dml.1.gz`.

    *Note: Ensure `$HOME/.local/bin` is in your PATH and `$HOME/.local/share/man` is in your MANPATH environment variable.*

## Usage

DML reads from standard input.

**Synopsis:**
```
dml [OPTIONS] < FILE
some_command | dml [OPTIONS]
```

**Options:**

*   `--colour COLOR`: Set the text colour for rendered LaTeX images. `COLOR` can be a named colour (e.g., "red", "blue") or a hex code (e.g., "#FF0000", "#0F0"). Defaults to "white".
*   `-c COLOR`: Short alias for `--colour`. If both are provided, `-c` takes precedence.
*   `--size SIZE`: Set the target terminal row height for rendered LaTeX images. `SIZE` is an integer. A value of `0` (default) uses 1 row for inline math and auto-sizes display math.
*   `-s SIZE`: Short alias for `--size`. If both are provided, `-s` takes precedence.
*   `--dpi DPI_VALUE`: Set the DPI (dots per inch) for rendering LaTeX images. `DPI_VALUE` is an integer. Defaults to `300`. Higher values produce sharper images but may be slower.
*   `-d DPI_VALUE`: Short alias for `--dpi`. If both are provided, `-d` takes precedence.
*   `--render-all-latex`: Render the entire input (including Markdown and text) as a single LaTeX document, which is then displayed as one image. This allows for consistent LaTeX font rendering throughout, but all text becomes part of an image.
*   `-l`: Short alias for `--render-all-latex`.
*   `--help` / `-h`: Displays help information about flags. (Standard Go flag behavior, prints to stderr).

**Examples:**

1.  **Display a file:**
    ```bash
    cat my_document.md | dml
    ```
    Or:
    ```bash
    dml < my_document.md
    ```

2.  **Pipe from `echo` with colored LaTeX:**
    ```bash
    echo 'This is **bold text** and inline math $E=mc^2$.' | dml --colour blue
    ```

3.  **Using the short color flag:**
    ```bash
    echo 'Display math: $$ \sum_{i=1}^{n} i = \frac{n(n+1)}{2} $$' | dml -c "#00FF00"
    ```

4.  **Render inline LaTeX image with a specific height (e.g., 2 terminal rows):**
    ```bash
    echo \'This is an inline formula $x^2$ sized to 2 rows.\' | dml -s 2
    ```

5.  **Render display math with a custom DPI (e.g., 150 DPI):**
    ```bash
    echo 'Display math at 150 DPI: $$ \int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2} $$' | dml --dpi 150
    ```

6.  **Render entire input as a single LaTeX image:**
    ```bash
    echo '# My Document\nThis is **bold text**, *italic text*, and some math $x^2 + y^2 = z^2$.\nAll of this will be one image.' | dml -l
    ```

7.  **Viewing the man page (after installation):**
    ```bash
    man dml
    ```

## Development

### Building from Source

To build DML from source:

```bash
# Clone the repository
git clone <repository_url>
cd dml

# Build the binary
./build.sh

# Run tests (if available)
go test ./...
```

## Author

Jamie Little
