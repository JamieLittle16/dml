# Regex Package

This package defines regular expression patterns used throughout DML for text processing and pattern matching.

## Key Components

- `patterns.go`: Contains regex pattern definitions for identifying and extracting various syntax elements

## Regex Patterns

The package provides several critical regex patterns:

### Math Expression Patterns

- `InlineMath`: Matches inline math expressions delimited by single dollar signs (`$...$`)
- `DisplayMath`: Matches display math expressions delimited by double dollar signs (`$$...$$`)
- `InlineMathParen`: Matches inline math expressions using parentheses (`\(...\)`)
- `DisplayMathBracket`: Matches display math expressions using brackets (`\[...\]`)

### Streaming Processing Patterns

- `StartDisplayMath`: Identifies the beginning of display math blocks during line-by-line processing
- `EndDisplayMath`: Identifies the end of display math blocks during line-by-line processing

These patterns use the `(?s)` flag where appropriate to ensure that dot (`.`) matches newlines, which is essential for multi-line math expressions.

## Usage

The regex patterns are used primarily in two contexts:

1. **Full Document Processing**: To identify and extract math expressions for rendering
2. **Streaming Processing**: To manage state during line-by-line processing of input

The streaming patterns are particularly important as they help maintain state when processing multi-line display math expressions, allowing DML to buffer content until it encounters a closing delimiter.

## Design Considerations

The patterns are carefully designed to:
- Handle edge cases like math expressions at line boundaries
- Support multiple delimiter styles (dollar signs, parentheses, and brackets)
- Balance performance with correctness for real-time processing

These regular expressions are a core part of DML's text processing pipeline, enabling accurate identification of math expressions embedded in regular text.