# Colour Package

This package provides colour management functionality for DML's LaTeX rendering.

## Key Components

- `ToHex()`: Converts colour names or hex codes to standardised hex format
- `ComplementHex()`: Calculates the complementary colour for a given hex colour
- `LaTeXColourDef()`: Generates LaTeX colour definitions for document templates
- `IsHexColour()`: Validates if a string is a properly formatted hex colour code

## Colour Processing

The package handles two types of colour specifications:
1. Named colours (e.g., "red", "blue", "green")
2. Hex codes (e.g., "#FF0000", "#00F", "#00FFFF")

Colour specifications are processed to ensure consistent formatting and compatibility with LaTeX. The package maintains a mapping of common colour names to their hex equivalents for convenience.

## Usage Example

```go
// Convert a colour name to hex
hexColour := colour.ToHex("blue")  // Returns "#0000FF"

// Get complementary colour
complement := colour.ComplementHex("#FF0000")  // Returns "#00FFFF"

// Generate LaTeX colour definition
colourDef := colour.LaTeXColourDef("textcolour", "#FF0000")
// Returns "\definecolour{textcolour}{HTML}{FF0000}"
```

This package is primarily used by the LaTeX rendering components to ensure consistent colour handling across different rendering modes.
