# Markdown Package

This package provides Markdown parsing, rendering, and processing functionality for DML.

## Key Components

- `markdown.go`: Contains the core Markdown processing functionality
  - AST (Abstract Syntax Tree) traversal
  - Terminal formatting with ANSI escape codes
  - LaTeX generation from Markdown content

## Functionality

### Markdown to LaTeX Conversion

The `GenerateLatexFromAST()` function traverses a Markdown AST and generates equivalent LaTeX syntax:
- Bold text (`**text**`) becomes `\textbf{text}`
- Italic text (`*text*`) becomes `\textit{text}`
- Code blocks become `\begin{verbatim}...\end{verbatim}`
- Headings are converted to appropriate LaTeX section commands
- Math expressions are preserved and properly formatted

This conversion is essential for the full document rendering mode, allowing mixed Markdown and LaTeX content to be rendered as a single cohesive document.

### Terminal Formatting

The `RenderMarkdownAST()` function applies ANSI terminal formatting:
- Bold text using ANSI code `\x1b[1m` (and `\x1b[22m` to reset)
- Italic text using ANSI code `\x1b[3m` (and `\x1b[23m` to reset)
- Other Markdown elements are processed recursively

### Markdown Processing Pipeline

The `ApplyFormatting()` function manages the overall Markdown processing:
1. Parses the Markdown text into an AST
2. Applies appropriate formatting transformations
3. Handles edge cases and cleanup for terminal display
4. Returns the formatted string ready for display

## Integration

This package works with other DML components:
- Uses standard Go Markdown parsing libraries
- Integrates with the LaTeX package for document generation
- Works alongside the terminal display components for consistent output

The Markdown processing pipeline is a core part of DML's functionality, enabling rich text display in terminal environments.