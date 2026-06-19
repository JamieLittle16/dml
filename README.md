# DML - Display Markdown & LaTeX

DML is a command-line tool that reads Markdown and LaTeX from standard input, renders math expressions as terminal images, and applies Markdown formatting to the output. It supports inline math (`$...$`), display math (`$$...$$`), formatted text (bold, italic, strikethrough), lists, blockquotes, tables with Unicode box-drawing, and links. LaTeX math is rendered as images using the Kitty terminal graphics protocol, with baseline-aligned inline math that stays on a single line.

## Code Structure

The codebase is organised in a modular structure:

- `cmd/dml/` - Main application entry point and CLI processing
- `internal/` - Core implementation packages:
  - `cache/` - Disk-backed LRU PNG cache (`~/.cache/dml/`)
  - `colour/` - Colour processing and management
  - `latex/` - LaTeX rendering, ImageMagick conversion, and cache integration
  - `markdown/` - Markdown processing, AST traversal, and table rendering
  - `regex/` - Regular expression patterns for math delimiter detection
  - `terminal/` - Terminal output, Kitty protocol, cell size queries, adaptive DPI
  - `unicode/` - Unicode fast-path rendering for simple math expressions

## Features

*   **Inline & display LaTeX math**: Renders `$formula$` (inline) and `$$formula$$` (display) as terminal images
*   **Baseline alignment**: Inline math is rendered to exact pixel baseline and stays on a single terminal line
*   **Markdown formatting**:
    - Bold: `**text**` or `__text__`
    - Italic: `*text*` or `_text_`
    - Strikethrough: `~~text~~`
    - Lists: `- item` or `1. item` (rendered with `•` or numbered)
    - Blockquotes: `> text` (rendered with `│` prefix)
    - Links: `[text](url)` (underlined with URL shown)
    - Inline code: `` `code` `` (reverse video)
    - Horizontal rules: `---` or `***` (renders as box-drawing line)
*   **Tables**: Markdown tables render with Unicode box-drawing characters
    ```
    ┌──────┬──────┐
    │ Name │ Type │
    ├──────┼──────┤
    │ foo  │ bar  │
    └──────┴──────┘
    ```
*   **Caching**: LRU disk cache at `~/.cache/dml/` (default 100 MB) for fast repeated renders
*   **Adaptive DPI**: Automatically scales render resolution to match terminal cell height (default range 96–600 DPI)
*   **Unicode fast path**: Simple expressions render as Unicode (e.g., `\alpha` → α, `x^2` → x²) without LaTeX pipeline
*   **Customisable text colour**: Set colour for LaTeX images with `--colour`
*   **Streams efficiently**: Processes input line-by-line with block buffering for complex Markdown structures

## Prerequisites

Before using DML, you need the following installed on your system:

1.  **Go**: Version 1.18 or higher (to build the tool).
2.  **pdflatex** and **latex** (dvipng): For compiling LaTeX expressions. These are part of standard TeX distributions like TeX Live or MiKTeX.
3.  **convert**: For converting PDF/DVI output into PNG images. This is part of the ImageMagick suite.
4.  **A Kitty-compatible terminal**: Required to display inline images. Ghostty and iTerm2 also support the Kitty graphics protocol.

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

    *Note: This requires sudo access. After installation, you might need to run `sudo mandb` (or `rehash` in some shells) for the system to recognise the new man page and command.*

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

*   `--colour COLOUR`: Set the text colour for rendered LaTeX images. `COLOUR` can be a named colour (e.g., "red", "blue") or a hex code (e.g., "#FF0000", "#0F0"). Defaults to "white".
*   `-c COLOUR`: Short alias for `--colour`. If both are provided, `-c` takes precedence.
*   `--dpi DPI_VALUE`: Set the DPI (dots per inch) for rendering LaTeX images. `DPI_VALUE` is an integer. Pass `0` (default) for adaptive DPI based on terminal cell height; otherwise specify a fixed DPI (96–600).
*   `-d DPI_VALUE`: Short alias for `--dpi`. If both are provided, `-d` takes precedence.
*   `--no-unicode`: Disable Unicode fast-path rendering; all math goes through LaTeX pipeline.
*   `--render-all-latex`: Render the entire input (including Markdown and text) as a single LaTeX document, which is then displayed as one image. This allows for consistent LaTeX font rendering throughout, but all text becomes part of an image.
*   `-l`: Short alias for `--render-all-latex`.
*   `--cache-stats`: Print cache statistics (hits, misses, size) and exit.
*   `--cache-clear`: Clear the render cache and exit.
*   `--cache-max-mb SIZE`: Set maximum cache size in MB (default 100). Cache uses LRU eviction when exceeded.
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

2.  **Pipe from `echo` with coloured LaTeX:**
    ```bash
    echo 'This is **bold text** and inline math $E=mc^2$.' | dml --colour blue
    ```

3.  **Using the short colour flag:**
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

## Caching

DML maintains a persistent disk cache at `~/.cache/dml/` to avoid re-rendering identical math expressions:

- **Cache key**: SHA-256 hash of LaTeX source + colour + DPI + display/inline + fuzz level
- **Storage**: PNG image + JSON metadata (dimensions, baseline offset, timestamps)
- **Eviction**: LRU (least-recently-used) when total PNG size exceeds limit
- **Default limit**: 100 MB
- **Non-fatal**: Cache failures do not break rendering; they're logged to stderr with `--debug`

**Cache management commands:**
```bash
dml --cache-stats              # Show cache statistics
dml --cache-clear             # Clear all cached entries
dml --cache-max-mb 200 < file # Set max cache size to 200 MB
```

## Development

### Building from Source

To build DML from source:

```bash
# Clone the repository
git clone https://github.com/JamieLittle16/dml
cd dml

# Build the binary
./build.sh

# Run tests (if available)
go test ./...
```

## Author

Jamie Little
