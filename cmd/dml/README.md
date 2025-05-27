# DML Command

This directory contains the main executable for the DML (Display Markdown & LaTeX) tool.

## Contents

- `main.go` - The main entry point for the DML application
  - Handles command-line argument parsing
  - Manages the processing pipeline
  - Coordinates between streaming and full document rendering modes
  - Implements error handling and user feedback

## Key Functions

- `main()`: Entry point that parses command-line flags and directs processing
- `processFullDocument()`: Handles rendering an entire document as a single LaTeX image
- `processStreamingDocument()`: Processes input line-by-line with state management for multi-line math
- `processInlineMath()`: Handles inline LaTeX math expressions within text

## Build Instructions

The main executable is built using the build script in the project root:

```bash
# From project root
./build.sh
```

This produces the `dml` binary which can be installed system-wide or used directly.