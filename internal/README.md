# Internal Packages

This directory contains the core implementation packages for DML that are used internally by the command-line interface.

## Packages

- `color/` - Color management functionality for LaTeX rendering
  - Handles color name to hex conversion
  - Provides complementary color calculation
  - Formats color definitions for LaTeX documents

- `latex/` - LaTeX document rendering and processing
  - Contains LaTeX document templates
  - Manages LaTeX document generation and compilation
  - Handles conversion of compiled documents to images
  - Provides text escaping utilities for LaTeX

- `markdown/` - Markdown text formatting and processing
  - Converts Markdown syntax to LaTeX for full document rendering
  - Applies ANSI styling for terminal display of Markdown formatting
  - Traverses Markdown AST (Abstract Syntax Tree) structures

- `regex/` - Regular expression patterns for text processing
  - Defines patterns for matching inline and display math expressions
  - Provides patterns for streaming text processing

- `terminal/` - Terminal-specific functionality
  - Implements Kitty terminal graphics protocol for image display
  - Manages terminal display characteristics

These packages are designed to be used together by the main application but maintain separation of concerns for better maintainability and testing.