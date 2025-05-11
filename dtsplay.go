package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/BourgeoisBear/rasterm"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

var (
    inlineMath  = regexp.MustCompile(`\$(.+?)\$`)
    	displayMath = regexp.MustCompile(`\$\$(.+?)\$\$`)
    )

    func renderMath(latex string, color string) ([]byte, error) {
    	var colorCmd string
    	if strings.HasPrefix(color, "#") {
    		colorCmd = fmt.Sprintf(`\color[HTML]{%s}`, strings.TrimPrefix(color, "#"))
    	} else {
    		colorCmd = fmt.Sprintf(`\color{%s}`, color)
    	}

    	// 1. Write minimal LaTeX document
    	tex := fmt.Sprintf(`\documentclass[preview]{standalone}
    \usepackage{amsmath}
    \usepackage{xcolor}
    \begin{document}
    %s$%s$
    \end{document}`, colorCmd, latex)

    	dir, err := ioutil.TempDir("", "dtsplay")
    if err != nil {
        return nil, err
    }
    texFile := dir + "/eq.tex"
    if err := ioutil.WriteFile(texFile, []byte(tex), 0644); err != nil {
        os.RemoveAll(dir) // Clean up if tex file writing fails
        return nil, err
    }

    // 2. Compile to PDF
    var stdout, stderr bytes.Buffer
    cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory", dir, texFile)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        // If pdflatex fails, do not remove the temp directory so logs can be inspected.
        // The error message now includes the path to the temp directory.
        return nil, fmt.Errorf("pdflatex failed for '%s': %v\nLaTeX STDOUT:\n%s\nLaTeX STDERR:\n%s\nTemp dir: %s", latex, err, stdout.String(), stderr.String(), dir)
    }

    // 3. Convert PDF → PNG
    pdfFile := dir + "/eq.pdf"
    pngFile := dir + "/eq.png"
    stdout.Reset() // Clear stdout buffer for the convert command
    stderr.Reset() // Clear stderr buffer for the convert command
    cmd = exec.Command("convert", "-density", "300", pdfFile, pngFile)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        // If convert command itself returns an error, keep temp dir for inspection.
        return nil, fmt.Errorf("convert command failed for PDF '%s': %v\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s", pdfFile, err, stdout.String(), stderr.String(), dir)
    }

    // Check if PNG file was actually created by convert, even if it exited 0
    if _, statErr := os.Stat(pngFile); os.IsNotExist(statErr) {
        // Convert "succeeded" (exit 0) but didn't create the file. Keep temp dir.
        return nil, fmt.Errorf("convert command appeared to succeed but did not create PNG '%s'.\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s", pngFile, stdout.String(), stderr.String(), dir)
    }

    // PNG file exists, now try to read it.
    imgData, readFileErr := ioutil.ReadFile(pngFile)
    if readFileErr != nil {
        // If reading fails (e.g., permissions, or file is corrupt), keep temp dir.
        return nil, fmt.Errorf("failed to read PNG file '%s': %v\nTemp dir: %s", pngFile, readFileErr, dir)
    }

    os.RemoveAll(dir) // Clean up ONLY on full success (pdflatex, convert, png exists, png read)
    return imgData, nil
}

func kittyInline(img []byte) (string, error) {
    // Use rasterm to emit Kitty protocol via stdout
    var sb strings.Builder
    opts := rasterm.KittyImgOpts{
        // Available fields in v1.1.1:
        // SrcX, SrcY, SrcWidth, SrcHeight, CellOffsetX, CellOffsetY,
        // DstCols, DstRows, ZIndex, ImageId, ImageNo, PlacementId
        // TransferMode and UseWindowSize are not available in this version.
        // Add any desired valid options here.
    }
    // KittyCopyPNGInline expects an io.Reader for the image data.
    // We capture the output into a string builder instead of writing to os.Stdout directly.
    err := rasterm.KittyCopyPNGInline(&sb, bytes.NewReader(img), opts)
    if err != nil {
        return "", err
    }
    	return sb.String(), nil
    }

    // renderMarkdownAST recursively traverses the AST and builds a string with ANSI codes.
    func renderMarkdownAST(node ast.Node, sb *strings.Builder) {
    	if node == nil {
    		return
    	}

    	switch n := node.(type) {
    	case *ast.Text:
    		sb.Write(n.Literal)
    	case *ast.Emph: // Handles *italic* and _italic_
    		sb.WriteString("\x1b[3m") // ANSI Italic on
    		for _, child := range n.GetChildren() {
    			renderMarkdownAST(child, sb)
    		}
    		sb.WriteString("\x1b[23m") // ANSI Italic off
    	case *ast.Strong: // Handles **bold** and __bold__
    		sb.WriteString("\x1b[1m") // ANSI Bold on
    		for _, child := range n.GetChildren() {
    			renderMarkdownAST(child, sb)
    		}
    		sb.WriteString("\x1b[22m") // ANSI Bold off
    	// For common containers, just recurse on children.
    	// This ensures their text content (including any nested Emph/Strong) is processed.
    	case *ast.Document, *ast.Paragraph, *ast.List, *ast.ListItem, *ast.Link, *ast.Image,
    		*ast.Code, *ast.CodeBlock, *ast.BlockQuote, *ast.Heading,
    		*ast.HorizontalRule, *ast.HTMLBlock, *ast.HTMLSpan,
    		*ast.Table, *ast.TableCell, *ast.TableHeader, *ast.TableRow,
    		*ast.Math, *ast.MathBlock: // Include Math/MathBlock in case parser produces them
    		for _, child := range n.GetChildren() {
    			renderMarkdownAST(child, sb)
    		}
    	default:
    		// For any other unhandled node type, if it's a container, process its children.
    		// This is a general fallback. If specific rendering is needed for other
    		// types, they should be added as explicit cases.
    		// fmt.Fprintf(os.Stderr, "DEBUG: Unhandled AST Node type: %T\n", n) // Optional debug
    		children := n.GetChildren()
    		for _, child := range children {
    			renderMarkdownAST(child, sb)
    		}
    	}
    }

    // applyMarkdownFormatting parses a line of Markdown text and converts basic
    // styling (bold, italics) to ANSI escape codes.
    func applyMarkdownFormatting(line string) string {
    	// Standard parser, no special extensions enabled by default.
    	// Extensions like parser.MathJax could be enabled if the Markdown
    	// parser should explicitly create *ast.Math or *ast.MathBlock nodes.
    	p := parser.New()
    	docNode := p.Parse([]byte(line))

    	var sb strings.Builder
    	renderMarkdownAST(docNode, &sb)
    	return sb.String()
    }

    func main() {
	// Command-line flags
	colourFlag := flag.String("colour", "white", "Set LaTeX text colour (e.g., red, #00FF00).") // Note: default flag name is "colour"
	cFlag := flag.String("c", "", "Short alias for --colour. Overrides --colour if set.")
	flag.Parse()

	// Determine the effective color to use
	effectiveColor := *colourFlag
	// If -c was provided (i.e., it's not its default empty string), it takes precedence.
	if *cFlag != "" {
		effectiveColor = *cFlag
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		// Step 1: Process LaTeX for rendering to Kitty protocol strings
		// Process display-math blocks ($$...$$) first
		line = displayMath.ReplaceAllStringFunc(line, func(match string) string {
			content := match[2 : len(match)-2] // Extract LaTeX content for display math
			img, err := renderMath(content, effectiveColor)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering display math '%s': %v\n", content, err)
				return match // Return original LaTeX string if rendering fails
			}
			kittyStr, err := kittyInline(img)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for display math '%s': %v\n", content, err)
				return match // Return original LaTeX string if Kitty protocol generation fails
			}
			return kittyStr // Replace LaTeX with Kitty image protocol string
		})

		// Process inline-math snippets ($...$) next
		line = inlineMath.ReplaceAllStringFunc(line, func(match string) string {
			content := match[1 : len(match)-1] // Extract LaTeX content for inline math
			img, err := renderMath(content, effectiveColor)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering inline math '%s': %v\n", content, err)
				return match // Return original LaTeX string if rendering fails
			}
			kittyStr, err := kittyInline(img)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math '%s': %v\n", content, err)
				return match // Return original LaTeX string if Kitty protocol generation fails
			}
			return kittyStr // Replace LaTeX with Kitty image protocol string
		})

		// Step 2: Apply Markdown formatting (italics, bold) to the line,
		// which now contains text and Kitty protocol strings from rendered LaTeX.
		line = applyMarkdownFormatting(line)

		// Print the fully processed line.
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error from scanner: %v\n", err)
	}
}
