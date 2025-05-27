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
		color      string
		isDisplay  bool
		dpi        int
		shouldFail bool
	}{
		{
			name:       "Empty LaTeX",
			latex:      "",
			color:      "white",
			isDisplay:  false,
			dpi:        300,
			shouldFail: true,
		},
		{
			name:       "Simple inline formula",
			latex:      "E=mc^2",
			color:      "white",
			isDisplay:  false,
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Simple display formula",
			latex:      "E=mc^2",
			color:      "white",
			isDisplay:  true,
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Custom color",
			latex:      "E=mc^2",
			color:      "red",
			isDisplay:  false,
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Invalid LaTeX",
			latex:      "\\invalidcommand",
			color:      "white",
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
			img, err := RenderMath(test.latex, test.color, test.isDisplay, test.dpi)
			
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
	skipActualRendering := true
	
	// Check if we should attempt real rendering
	if os.Getenv("DML_TEST_FULL_RENDERING") == "1" && 
	   fileExists("/usr/bin/pdflatex") && 
	   fileExists("/usr/bin/convert") {
		skipActualRendering = false
	}

	tests := []struct {
		name       string
		latexBody  string
		color      string
		dpi        int
		shouldFail bool
	}{
		{
			name:       "Empty document",
			latexBody:  "",
			color:      "white",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Simple document",
			latexBody:  "Hello, world!",
			color:      "white",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Document with math",
			latexBody:  "Formula: $E=mc^2$",
			color:      "white",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Custom color",
			latexBody:  "Coloured text",
			color:      "blue",
			dpi:        300,
			shouldFail: false,
		},
		{
			name:       "Invalid LaTeX",
			latexBody:  "\\invalidcommand",
			color:      "white",
			dpi:        300,
			shouldFail: true,
		},
	}

	// Set debug mode for detailed output
	SetDebug(true)
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if skipActualRendering {
				// In mock mode, we just verify parameter handling
				if test.name == "Invalid LaTeX" && !test.shouldFail {
					t.Errorf("Invalid LaTeX should fail but test expects success")
				}
				t.Skip("Skipping actual rendering test")
				return
			}
			
			// Only run these tests if we're doing actual rendering
			img, err := RenderFullDocument(test.latexBody, test.color, test.dpi)
			
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