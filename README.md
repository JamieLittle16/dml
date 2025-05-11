# DML - Display Markdown & LaTeX

DML is a command-line tool that reads text from standard input, processes basic Markdown (bold/italic) and LaTeX math expressions (both inline and display), and prints the result to standard output. LaTeX math is rendered as images using the Kitty terminal graphics protocol.

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
    git clone <repository_url>
    cd dml
    ```
    If you just have the source files, navigate to the `Projects/dml` directory.

2.  **Build the Binary**:
    From the `Projects/dml` directory:
    ```bash
    go build -o dml dtsplay.go
    ```
    This will create an executable file named `dml` in the current directory.

3.  **Install the Binary and Man Page (System-Wide)**:
    To make `dml` available globally and install its man page, run the following commands from the `Projects/dml` directory (you will likely need `sudo`):

    ```bash
    sudo cp dml /usr/local/bin/dml
    sudo mkdir -p /usr/local/share/man/man1
    sudo cp man/dml.1 /usr/local/share/man/man1/dml.1
    sudo gzip /usr/local/share/man/man1/dml.1
    ```
    *Note: If `/usr/local/share/man/man1` does not exist, the `mkdir -p` command will create it.*
    *After installation, you might need to run `sudo mandb` (or `rehash` in some shells) for the system to recognize the new man page and command.*

4.  **User-Specific Installation (Optional)**:
    Alternatively, you can install it to a user-specific directory like `$HOME/.local/bin` (ensure this is in your PATH):
    ```bash
    mkdir -p $HOME/.local/bin
    cp dml $HOME/.local/bin/dml
    # For the man page:
    mkdir -p $HOME/.local/share/man/man1
    cp man/dml.1 $HOME/.local/share/man/man1/dml.1
    gzip $HOME/.local/share/man/man1/dml.1
    # Ensure $HOME/.local/share/man is in your MANPATH environment variable.
    ```

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

6.  **Viewing the man page (after installation):**
    ```bash
    man dml
    ```

## Author

Jamie Little
