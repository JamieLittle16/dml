# Color Package

This package provides color management functionality for DML's LaTeX rendering.

## Key Components

- `ToHex()`: Converts color names or hex codes to standardized hex format
- `ComplementHex()`: Calculates the complementary color for a given hex color
- `LaTeXColorDef()`: Generates LaTeX color definitions for document templates
- `IsHexColor()`: Validates if a string is a properly formatted hex color code

## Color Processing

The package handles two types of color specifications:
1. Named colors (e.g., "red", "blue", "green")
2. Hex codes (e.g., "#FF0000", "#00F", "#00FFFF")

Color specifications are processed to ensure consistent formatting and compatibility with LaTeX. The package maintains a mapping of common color names to their hex equivalents for convenience.

## Usage Example

```go
// Convert a color name to hex
hexColor := color.ToHex("blue")  // Returns "#0000FF"

// Get complementary color
complement := color.ComplementHex("#FF0000")  // Returns "#00FFFF"

// Generate LaTeX color definition
colorDef := color.LaTeXColorDef("textcolor", "#FF0000")
// Returns "\definecolor{textcolor}{HTML}{FF0000}"
```

This package is primarily used by the LaTeX rendering components to ensure consistent color handling across different rendering modes.