# LaTeX Package

This package handles LaTeX document generation, rendering, and processing for DML.

## Key Components

- `template.go`: Contains LaTeX document templates for both inline math and full document rendering
- `render.go`: Implements the core rendering functionality that converts LaTeX to images
- `escape.go`: Provides utilities for escaping special characters in LaTeX

## Functionality

### Templates

The package provides two main templates:
- A minimal template for rendering individual math expressions
- A comprehensive template for rendering full documents with proper structure

### Rendering

The rendering process follows these steps:
1. Template selection and content preparation
2. LaTeX document generation with proper color settings
3. Compilation using `pdflatex` to create a PDF
4. Conversion to PNG format using ImageMagick's `convert` utility
5. Image processing for transparency and proper display

The package offers two main rendering functions:
- `RenderMath()`: For individual math expressions (inline or display)
- `RenderFullDocument()`: For entire documents with mixed content

### Escaping

LaTeX has numerous special characters that need escaping when used in regular text. The escaping functionality ensures that text content is properly formatted for LaTeX compilation by handling characters like:
- Backslashes, braces, and brackets
- Dollar signs, percent signs, and ampersands
- Other special characters that have meaning in LaTeX

## Debug Support

The package includes comprehensive debug logging to help diagnose rendering issues:
- LaTeX document generation and content
- Temporary file management
- Compilation output capture
- Image conversion process monitoring

This information is critical for troubleshooting rendering problems, especially when running on systems with different LaTeX installations.