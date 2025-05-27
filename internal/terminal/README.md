# Terminal Package

This package provides terminal-specific functionality for DML, focusing on displaying rendered LaTeX images in the terminal.

## Key Components

- `kitty.go`: Implements the Kitty terminal graphics protocol for displaying images inline with text

## Functionality

### Kitty Graphics Protocol

The terminal package implements the [Kitty Graphics Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/) which allows displaying raster images directly in the terminal:

- `KittyInline()`: Generates the Kitty protocol escape sequences for displaying an image
  - Handles PNG image data from LaTeX rendering
  - Configures proper sizing for both inline and display math
  - Manages terminal-specific formatting like newlines and escape characters

### Image Display Configuration

The package provides careful handling of different math display modes:

- **Inline Math**: Typically sized to a single terminal row and integrated with surrounding text
- **Display Math**: Automatically sized based on content with proper vertical spacing

### Debug Support

The terminal package includes debug functionality that logs:
- Image protocol generation parameters
- Terminal row sizing decisions
- Output formatting details

This information helps diagnose display issues across different terminal types and configurations.

## Integration

This package is the final step in the DML rendering pipeline:
1. LaTeX expressions are identified in the input text
2. These expressions are rendered to PDF and then PNG format
3. The terminal package takes these PNG images and generates the appropriate escape sequences
4. These sequences are written to the terminal, displaying the images inline with text

The terminal output is carefully managed to ensure proper text flow around both inline and display math expressions.