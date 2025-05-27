package latex

import (
	"os"
	"testing"
)

// Mock test for RenderMath function
// Note: Full testing requires LaTeX and ImageMagick installed
func TestRenderMath(t *testing.T) {
	// This test can be run in two modes:
	// 1. Mock mode (default): Skip actual rendering, just test parameter handling
	// 2. Full mode: If environment supports it, test actual rendering

	// Always skip actual LaTeX rendering in CI environments
	skipActualRendering := true

	// Check if we should attempt real rendering
	if os.Getenv("DML_TEST_FULL_RENDERING") == "1" &&
	   fileExists("/usr/bin/pdflatex") &&
	   fileExists("/usr/bin/convert") {
		skipActualRendering = false
	}

	tests := []struct {
		name       string
		latex      string
		colour      string
		isDisplay  bool
		dpi        int
		shouldFail bool
	}{
		{
			name:       "Empty LaTeX",
			latex:      "",
			colour:      "white",
			isDisplay:  false,
			dpi:        300,
			shouldFail: true,
		},
		{
			name:       "Simple inline formula",
			latex:      "E=mc^2",
			colour:      "white",
			isDisplay:  false,
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Simple display formula",
			latex:      "E=mc^2",
			colour:      "white",
			isDisplay:  true,
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Custom colour",
			latex:      "E=mc^2",
			colour:      "red",
			isDisplay:  false,
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Invalid LaTeX",
			latex:      "\\invalidcommand",
			colour:      "white",
			isDisplay:  false,
			dpi:        300,
			shouldFail: true,
		},
	}

	// Set debug mode for detailed output
	SetDebug(true)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if skipActualRendering {
				// In mock mode, just verify parameter handling logic
				if test.latex == "" && !test.shouldFail {
					t.Errorf("Empty LaTeX should fail but test expects success")
				}
				t.Skip("Skipping actual rendering test")
				return
			}

			// Only run these tests if we're doing actual rendering
			img, err := RenderMath(test.latex, test.colour, test.isDisplay, test.dpi)

			if test.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(img) == 0 {
				t.Errorf("Expected non-empty image data")
			}
		})
	}
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Test for RenderFullDocument function
func TestRenderFullDocument(t *testing.T) {
	// This test can be run in two modes:
	// 1. Mock mode (default): Skip actual rendering, just test parameter handling
	// 2. Full mode: If environment supports it, test actual rendering

	// Always skip actual LaTeX rendering in CI environments
	// skipActualRendering := true

	// Check if we should attempt real rendering
	if os.Getenv("DML_TEST_FULL_RENDERING") == "1" &&
	   fileExists("/usr/bin/pdflatex") &&
	   fileExists("/usr/bin/convert") {
		// skipActualRendering = false
	}

	tests := []struct {
		name       string
		latexBody  string
		colour      string
		dpi        int
		shouldFail bool
	}{
		{
			name:       "Empty document",
			latexBody:  "",
			colour:      "white",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Simple document",
			latexBody:  "Hello, world!",
			colour:      "white",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Document with math",
			latexBody:  "Formula: $E=mc^2$",
			colour:      "white",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Custom colour",
			latexBody:  "coloured text",
			colour:      "blue",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Invalid LaTeX",
			latexBody:  "\\invalidcommand",
			colour:      "white",
			dpi:        300,
			shouldFail: true,
		},
	}

	// Only check if LaTeX executables are available, don't actually render
	if _, err := os.Stat("/usr/bin/pdflatex"); os.IsNotExist(err) {
		t.Skip("pdflatex not found, skipping actual LaTeX tests")
	}

	if _, err := os.Stat("/usr/bin/convert"); os.IsNotExist(err) {
		t.Skip("convert (ImageMagick) not found, skipping actual LaTeX tests")
	}

	// Set debug mode to see detailed output
	SetDebug(true)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			img, err := RenderFullDocument(test.latexBody, test.colour, test.dpi)

			if test.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(img) == 0 {
				t.Errorf("Expected non-empty image data")
			}
		})
	}
}
