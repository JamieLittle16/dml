// Package latex provides LaTeX rendering functionality for DML
package latex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"dml/internal/colour"
)

var isDebug bool

// SetDebug enables or disables debug mode
func SetDebug(debug bool) {
	isDebug = debug
}

// RenderMath renders a LaTeX math expression to a PNG image
func RenderMath(latex string, colourStr string, isDisplay bool, dpi int, fuzzLevel string) ([]byte, error) {
	// Skip empty latex content
	latex = strings.TrimSpace(latex)
	if latex == "" {
		return nil, fmt.Errorf("empty LaTeX content")
	}

	// Add a small amount of spacing around display math for better rendering
	if isDisplay {
		latex = " " + latex + " "
	}

	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: renderMath called with isDisplay=%v, dpi=%d\n", isDisplay, dpi)
	}

	// Set default colour and prepare background colour
	if colourStr == "" {
		colourStr = "white"
	}
	bg := "black"
	transparent := "black"
	latexcolourDefs := ""

	// Process colours
	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: Processing colour: '%s'\n", colourStr)
	}

	hexcolour := colour.ToHex(colourStr)
	if hexcolour != "" {
		comp := colour.ComplementHex(hexcolour)
		latexcolourDefs = colour.LaTeXcolourDef("usercolour", hexcolour) + colour.LaTeXcolourDef("bgcolour", comp)
		colourStr = "usercolour"
		bg = "bgcolour"
		transparent = comp

		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: colour converted to hex: '%s', complement: '%s'\n", hexcolour, comp)
		}
	} else {
		// fallback: use white text on black bg
		colourStr = "white"
		bg = "black"
		transparent = "black"

		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: Using fallback colours: text='%s', bg='%s'\n", colourStr, bg)
		}
	}

	// Prepare LaTeX content
	var mathContent string
	if isDisplay {
		mathContent = fmt.Sprintf(`\[%s\]`, latex)
		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: Creating display math content with \\[ ... \\]\n")
		}
	} else {
		mathContent = fmt.Sprintf(`$%s$`, latex)
		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: Creating inline math content with $ ... $\n")
		}
	}

	// Fill the template
	tex := fmt.Sprintf(TexTemplate, latexcolourDefs, bg, colourStr, mathContent)

	// Create temporary directory
	dir, err := ioutil.TempDir("", "dml")
	if err != nil {
		return nil, err
	}

	texFile := dir + "/eq.tex"
	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: Created temp directory for LaTeX rendering: %s\n", dir)
	}

	// Write the TeX file
	if err := ioutil.WriteFile(texFile, []byte(tex), 0644); err != nil {
		os.RemoveAll(dir) // Clean up if tex file writing fails
		return nil, err
	}

	// Compile to PDF
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory", dir, texFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// If pdflatex fails, do not remove the temp directory so logs can be inspected
		return nil, fmt.Errorf("pdflatex failed for '%s': %v\nLaTeX STDOUT:\n%s\nLaTeX STDERR:\n%s\nTemp dir: %s",
			latex, err, stdout.String(), stderr.String(), dir)
	}

	// Convert PDF to PNG
	pdfFile := dir + "/eq.pdf"
	pngFile := dir + "/eq.png"
	stdout.Reset()
	stderr.Reset()

	// Determine fuzz level to use
	effectiveFuzz := fuzzLevel
	if effectiveFuzz == "" {
		effectiveFuzz = "1%" // Default fuzz level
	}

	// Convert with extremely aggressive background removal
	cmd = exec.Command("convert",
		"-density", fmt.Sprintf("%d", dpi),
		"-alpha", "on",
		"-background", "none",
		"-trim",
		"+repage",
		"-fuzz", effectiveFuzz,
		"-transparent", transparent,
		"-fuzz", effectiveFuzz,
		"-transparent", transparent,
		"-channel", "Alpha",
		"-threshold", "50%",
		"+channel",
		"-background", "none",
		"-flatten",
		pdfFile, pngFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("convert command failed for PDF '%s': %v\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s",
			pdfFile, err, stdout.String(), stderr.String(), dir)
	}

	// Check if PNG file was created
	if _, statErr := os.Stat(pngFile); os.IsNotExist(statErr) {
		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: PNG file not found after conversion: %s\n", pngFile)
		}
		return nil, fmt.Errorf("convert command appeared to succeed but did not create PNG '%s'.\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s\nStat error: %v",
			pngFile, stdout.String(), stderr.String(), dir, statErr)
	}

	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: PNG file successfully created: %s\n", pngFile)
	}

	// Read the PNG file
	imgData, readFileErr := ioutil.ReadFile(pngFile)
	if readFileErr != nil {
		return nil, fmt.Errorf("failed to read PNG file '%s': %v\nTemp dir: %s", pngFile, readFileErr, dir)
	}

	os.RemoveAll(dir) // Clean up only on full success
	return imgData, nil
}

// RenderFullDocument renders an entire document as a single LaTeX image
func RenderFullDocument(latexBody string, colourStr string, dpi int, fuzzLevel string) ([]byte, error) {
	if colourStr == "" {
		colourStr = "white"
	}
	bg := "black"
	transparent := "black"
	latexcolourDefs := ""

	// Process colours
	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: Processing colour for full document: '%s'\n", colourStr)
	}

	hexcolour := colour.ToHex(colourStr)
	if hexcolour != "" {
		comp := colour.ComplementHex(hexcolour)
		latexcolourDefs = colour.LaTeXcolourDef("usercolour", hexcolour) + colour.LaTeXcolourDef("bgcolour", comp)
		colourStr = "usercolour"
		bg = "bgcolour"
		transparent = comp

		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: Full doc colour converted to hex: '%s', complement: '%s'\n", hexcolour, comp)
		}
	} else {
		// fallback: use white text on black bg
		colourStr = "white"
		bg = "black"
		transparent = "black"

		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: Using fallback colours for full doc: text='%s', bg='%s'\n", colourStr, bg)
		}
	}

	// Fill the template
	tex := fmt.Sprintf(FullDocTemplate, latexcolourDefs, bg, colourStr, latexBody)

	// Create temporary directory
	dir, err := ioutil.TempDir("", "dml-full")
	if err != nil {
		return nil, err
	}

	texFile := dir + "/fulldoc.tex"
	if err := ioutil.WriteFile(texFile, []byte(tex), 0644); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}

	// Compile to PDF
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory", dir, texFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pdflatex failed for full document: %v\nLaTeX STDOUT:\n%s\nLaTeX STDERR:\n%s\nTemp dir: %s",
			err, stdout.String(), stderr.String(), dir)
	}

	// Convert PDF to PNG
	pdfFile := dir + "/fulldoc.pdf"
	pngFile := dir + "/fulldoc.png"
	stdout.Reset()
	stderr.Reset()

	// Determine fuzz level to use
	effectiveFuzz := fuzzLevel
	if effectiveFuzz == "" {
		effectiveFuzz = "1%" // Default fuzz level
	}

	cmd = exec.Command("convert",
		"-density", fmt.Sprintf("%d", dpi),
		"-quality", "100",
		"-trim",
		"+repage",
		"-fuzz", effectiveFuzz,
		"-transparent", transparent,
		"-fuzz", effectiveFuzz,
		"-transparent", transparent,
		"-channel", "Alpha",
		"-threshold", "50%",
		"+channel",
		"-background", "none",
		"-flatten",
		pdfFile, pngFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("convert command failed for PDF '%s': %v\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s",
			pdfFile, err, stdout.String(), stderr.String(), dir)
	}

	// Check if PNG file exists
	if _, statErr := os.Stat(pngFile); os.IsNotExist(statErr) {
		return nil, fmt.Errorf("convert command appeared to succeed but did not create PNG '%s'.\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s",
			pngFile, stdout.String(), stderr.String(), dir)
	}

	// Read the PNG file
	imgData, readFileErr := ioutil.ReadFile(pngFile)
	if readFileErr != nil {
		return nil, fmt.Errorf("failed to read PNG file '%s': %v\nTemp dir: %s", pngFile, readFileErr, dir)
	}

	os.RemoveAll(dir)
	return imgData, nil
}
