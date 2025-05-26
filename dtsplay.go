package main

import (
	"bufio" // Added for buffered input/output streaming
	"bytes"
	"flag"
	"fmt"
	"io" // Added for EOF handling
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/BourgeoisBear/rasterm"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// Returns true if s is a hex color (#RRGGBB or #RGB)
func isHexColor(s string) bool {
	matched, _ := regexp.MatchString(`^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$`, s)
	return matched
}

// Map of named colors to hex codes (CSS/X11 names, can expand as needed)
var namedColorHex = map[string]string{
	"black":   "#000000",
	"white":   "#FFFFFF",
	"red":     "#FF0000",
	"green":   "#00FF00",
	"blue":    "#0000FF",
	"yellow":  "#FFFF00",
	"cyan":    "#00FFFF",
	"magenta": "#FF00FF",
	"gray":    "#808080",
	"grey":    "#808080",
	"orange":  "#FFA500",
	"purple":  "#800080",
	"brown":   "#A52A2A",
	"pink":    "#FFC0CB",
	"lime":    "#00FF00",
	"navy":    "#000080",
	"teal":    "#008080",
	"maroon":  "#800000",
	"olive":   "#808000",
	"silver":  "#C0C0C0",
	// Add more as needed
}

// Returns a hex code for a color string (hex or named), or empty string if unknown
func colorToHex(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if isHexColor(s) {
		return expandHexColor(s)
	}
	if hex, ok := namedColorHex[s]; ok {
		return hex
	}
	return ""
}

// Expands #RGB to #RRGGBB
func expandHexColor(s string) string {
	if len(s) == 4 {
		return "#" + strings.Repeat(string(s[1]), 2) +
			strings.Repeat(string(s[2]), 2) +
			strings.Repeat(string(s[3]), 2)
	}
	return s
}

// Returns the complement of a #RRGGBB hex color as #RRGGBB
func complementHexColor(s string) string {
	s = expandHexColor(s)
	r, _ := strconv.ParseUint(s[1:3], 16, 8)
	g, _ := strconv.ParseUint(s[3:5], 16, 8)
	b, _ := strconv.ParseUint(s[5:7], 16, 8)
	return fmt.Sprintf("#%02X%02X%02X", 0xFF^r, 0xFF^g, 0xFF^b)
}

// Returns the color name and LaTeX color definition for a hex color
func latexColorDef(name, hex string) string {
	return fmt.Sprintf("\\definecolor{%s}{HTML}{%s}\n", name, hex[1:])
}

var (
	// Added (?s) to make '.' match newlines for multi-line LaTeX expressions.
	inlineMath         = regexp.MustCompile(`(?s)\$(.+?)\$`)
	displayMath        = regexp.MustCompile(`(?s)\$\$(.+?)\$\$`)
	inlineMathParen    = regexp.MustCompile(`(?s)\\\((.+?)\\\)`)    // For \( ... \)
	displayMathBracket = regexp.MustCompile(`(?s)\\\[(.+?)\\\]`) // For \[ ... \]
)

// LaTeX special characters that need escaping. Note: `\` is handled by replacing with `\textbackslash{}`.
// Order matters for some replacements (e.g., `\` before other chars that might use it).
var latexEscaper = strings.NewReplacer(
	`\`, `\textbackslash{}`, // Must be first
	`&`, `\&`,
	`%`, `\%`,
	`$`, `\$`,
	`#`, `\#`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
	`~`, `\textasciitilde{}`,
	`^`, `\textasciicircum{}`,
	`[`, `{[}`,
	`]`, `{]}`,
	`|`, `{\vert}`,
	`/`, `{/}`,
	`[`, `{[}`,
	`]`, `{]}`,
	`|`, `{|}`,
	`/`, `{/}`,
)

func escapeLatex(s string) string {
	return latexEscaper.Replace(s)
}

func generateLatexFromAST(node ast.Node, sb *strings.Builder) {
	if node == nil {
		return
	}

	// Helper to process children
	processChildren := func(n ast.Node) {
		for _, child := range n.GetChildren() {
			generateLatexFromAST(child, sb)
		}
	}

	switch n := node.(type) {
	case *ast.Document:
		processChildren(n)
	case *ast.Text:
		sb.WriteString(escapeLatex(string(n.Literal)))
	case *ast.Emph: // *italic* or _italic_
		sb.WriteString(`\textit{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Strong: // **bold** or __bold__
		sb.WriteString(`\textbf{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Math: // Inline math $...$ or \(...\) from MathJax extension
		sb.WriteString(`$`)
		sb.WriteString(string(n.Literal)) // n.Literal contains the raw math content
		sb.WriteString(`$`)
	case *ast.MathBlock: // Display math $$...$$ or \[...\] from MathJax extension
		sb.WriteString(`$$`)
		sb.WriteString(string(n.Literal)) // n.Literal contains the raw math content
		sb.WriteString(`$$`)
	case *ast.Paragraph:
		processChildren(n)
		sb.WriteString("\n\\par\n\n") // LaTeX paragraph break
	case *ast.Softbreak: // Typically a newline in Markdown source
		sb.WriteString("\\\\\n") // LaTeX line break
	case *ast.Hardbreak: // Explicit line break (e.g., two spaces at end of line)
		sb.WriteString("\\\\\n") // LaTeX line break
	case *ast.Code: // Inline code
		sb.WriteString(`\texttt{`)
		sb.WriteString(escapeLatex(string(n.Literal)))
		sb.WriteString(`}`)
	case *ast.CodeBlock: // Code block
		sb.WriteString("\n\\begin{verbatim}\n")
		sb.WriteString(string(n.Literal)) // No need to escape in verbatim
		sb.WriteString("\\end{verbatim}\n")
	case *ast.Heading: // Section headings
		level := n.Level
		prefix := ""
		switch level {
		case 1:
			prefix = "\\section*{"
		case 2:
			prefix = "\\subsection*{"
		case 3:
			prefix = "\\subsubsection*{"
		case 4, 5, 6:
			prefix = "\\paragraph*{"
		}
		sb.WriteString(prefix)
		processChildren(n)
		sb.WriteString("}\n")
	// Default handling for other common nodes: just process their children.
	case *ast.List, *ast.ListItem, *ast.Link, *ast.Image,
		*ast.BlockQuote, *ast.HorizontalRule, *ast.HTMLBlock, *ast.HTMLSpan,
		*ast.Table, *ast.TableCell, *ast.TableHeader, *ast.TableRow:
		processChildren(n) // For now, just try to render their content
	default:
		// For any other unhandled node type, if it's a container, process its children.
		children := n.GetChildren()
		for _, child := range children {
			generateLatexFromAST(child, sb)
		}
	}
}

func renderFullLatexDocument(latexBody string, color string, dpi int) ([]byte, error) {
	if color == "" {
		color = "white"
	}
	bg := "black"
	transparent := "black"
	latexColorDefs := ""

	hexColor := colorToHex(color)
	if hexColor != "" {
		comp := complementHexColor(hexColor)
		latexColorDefs = latexColorDef("usercolor", hexColor) + latexColorDef("bgcolor", comp)
		color = "usercolor"
		bg = "bgcolor"
		transparent = comp
	} else {
		// fallback: use white text on black bg
		color = "white"
		bg = "black"
		transparent = "black"
	}

	texTemplate := `\documentclass[border=3pt,preview]{standalone}
\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{amsfonts}
\usepackage{mathtools}
\usepackage[dvipsnames,svgnames,table]{xcolor}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{lmodern}
\usepackage{verbatim}
%s
\begin{document}
\pagecolor{%s}
\color{%s}
%s
\end{document}`
	tex := fmt.Sprintf(texTemplate, latexColorDefs, bg, color, latexBody)

	dir, err := ioutil.TempDir("", "dtsplay-full")
	if err != nil {
		return nil, err
	}
	texFile := dir + "/fulldoc.tex"
	if err := ioutil.WriteFile(texFile, []byte(tex), 0644); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory", dir, texFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pdflatex failed for full document: %v\nLaTeX STDOUT:\n%s\nLaTeX STDERR:\n%s\nTemp dir: %s", err, stdout.String(), stderr.String(), dir)
	}

	pdfFile := dir + "/fulldoc.pdf"
	pngFile := dir + "/fulldoc.png"
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("convert", 
		"-density", fmt.Sprintf("%d", dpi),
		"-quality", "100", 
		"-trim",             // Remove any excess whitespace
		"+repage",           // Reset the page after trimming
		"-transparent", transparent,
		pdfFile, pngFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("convert command failed for PDF '%s': %v\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s", pdfFile, err, stdout.String(), stderr.String(), dir)
	}

	if _, statErr := os.Stat(pngFile); os.IsNotExist(statErr) {
		return nil, fmt.Errorf("convert command appeared to succeed but did not create PNG '%s'.\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s", pngFile, stdout.String(), stderr.String(), dir)
	}

	imgData, readFileErr := ioutil.ReadFile(pngFile)
	if readFileErr != nil {
		return nil, fmt.Errorf("failed to read PNG file '%s': %v\nTemp dir: %s", pngFile, readFileErr, dir)
	}

	os.RemoveAll(dir)
	return imgData, nil
}

func renderMath(latex string, color string, isDisplay bool, dpi int) ([]byte, error) {
    // Skip empty latex content
    latex = strings.TrimSpace(latex)
    if latex == "" {
        return nil, fmt.Errorf("empty LaTeX content")
    }
    // Get debugging flag value from environment - needed for functions outside main
    debugEnv := os.Getenv("DML_DEBUG")
    isDebug := debugEnv == "1" || debugEnv == "true" || debugEnv == "yes"

    if isDebug {
        fmt.Fprintf(os.Stderr, "DEBUG: renderMath called with isDisplay=%v, dpi=%d\n", isDisplay, dpi)
    }
	if color == "" {
		color = "white"
	}
	bg := "black"
	transparent := "black"
	latexColorDefs := ""

	hexColor := colorToHex(color)
	if hexColor != "" {
		comp := complementHexColor(hexColor)
		latexColorDefs = latexColorDef("usercolor", hexColor) + latexColorDef("bgcolor", comp)
		color = "usercolor"
		bg = "bgcolor"
		transparent = comp
	} else {
		// fallback: use white text on black bg
		color = "white"
		bg = "black"
		transparent = "black"
	}

	texTemplate := `\documentclass[border=4pt,preview]{standalone}
\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{amsfonts}
\usepackage{mathtools}
\usepackage[dvipsnames,svgnames,table]{xcolor}
%s
\begin{document}
\pagecolor{%s}
\color{%s}
%s
\end{document}`
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
	tex := fmt.Sprintf(texTemplate, latexColorDefs, bg, color, mathContent)

	dir, err := ioutil.TempDir("", "dtsplay")
	if err != nil {
		return nil, err
	}
	texFile := dir + "/eq.tex"
	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: Created temp directory for LaTeX rendering: %s\n", dir)
	}
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

	// Enhanced convert command with better options for rendering quality
	cmd = exec.Command("convert", 
		"-density", fmt.Sprintf("%d", dpi),
		"-quality", "100",
		"-trim",             // Remove any excess whitespace
		"+repage",           // Reset the page after trimming
		"-bordercolor", transparent,
		"-border", "1x1",    // Add a tiny border for better spacing
		"-transparent", transparent,
		pdfFile, pngFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// If convert command itself returns an error, keep temp dir for inspection.
		return nil, fmt.Errorf("convert command failed for PDF '%s': %v\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s", pdfFile, err, stdout.String(), stderr.String(), dir)
	}

	// Check if PNG file was actually created by convert, even if it exited 0
	if _, statErr := os.Stat(pngFile); os.IsNotExist(statErr) {
		// Convert "succeeded" (exit 0) but didn't create the file. Keep temp dir.
		if isDebug {
			fmt.Fprintf(os.Stderr, "DEBUG: PNG file not found after conversion: %s\n", pngFile)
		}
		return nil, fmt.Errorf("convert command appeared to succeed but did not create PNG '%s'.\nConverter STDOUT:\n%s\nConverter STDERR:\n%s\nTemp dir: %s\nStat error: %v", pngFile, stdout.String(), stderr.String(), dir, statErr)
	}
	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: PNG file successfully created: %s\n", pngFile)
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

// kittyInline generates the Kitty graphics protocol string for the given image bytes.
// It carefully handles the output to ensure no unwanted trailing characters appear.
func kittyInline(img []byte, isDisplayMath bool, userTargetRows int) (string, error) {
    var sb strings.Builder
    opts := rasterm.KittyImgOpts{
        // Available fields in v1.1.1:
        // SrcX, SrcY, SrcWidth, SrcHeight, CellOffsetX, CellOffsetY,
        // DstCols, DstRows, ZIndex, ImageId, ImageNo, PlacementId
        // TransferMode and UseWindowSize are not available in this version.
    }

    // Get debug flag value from environment
    debugEnv := os.Getenv("DML_DEBUG")
    isDebug := debugEnv == "1" || debugEnv == "true" || debugEnv == "yes"

    actualDstRows := 0
    if userTargetRows > 0 {
        actualDstRows = userTargetRows
    } else { // User did not specify a size, use our defaults
        if isDisplayMath {
            // For display math, 0 might let rasterm/kitty decide size, or pick a larger default.
            // Let's use 0 for auto-sizing display math by default if not specified.
            actualDstRows = 0
        } else { // Inline math
            actualDstRows = 1 // Default to 1 row for inline math
        }
    }

    if actualDstRows > 0 {
        opts.DstRows = uint32(actualDstRows)
    }

	// rasterm.KittyCopyPNGInline expects an io.Reader for the image data.
	// We capture its output (the Kitty protocol string) into a strings.Builder.
	err := rasterm.KittyCopyPNGInline(&sb, bytes.NewReader(img), opts)
	if err != nil {
		return "", fmt.Errorf("rasterm.KittyCopyPNGInline failed: %v", err)
	}
	kittyStr := sb.String()

	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: Generated Kitty protocol with options: rows=%d, isDisplay=%v\n",
			actualDstRows, isDisplayMath)
	}

	// Handle newlines for display math and inline math differently
	if isDisplayMath {
		// For display math, ensure there's exactly one trailing newline
		kittyStr = strings.TrimRight(kittyStr, "\n") + "\n"
	} else {
		// For inline math, ensure there are no trailing newlines or special characters
		kittyStr = strings.TrimRight(kittyStr, "\n")
	}
	
	// Remove any unwanted characters that might be present in the Kitty protocol output
	kittyStr = strings.TrimSuffix(kittyStr, "%")
	kittyStr = strings.TrimSuffix(kittyStr, "\x00")
	
	return kittyStr, nil
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
    	// Standard parser, no special extensions enabled by default for this function.
    	// This function is only used in the non-full-LaTeX mode where LaTeX
    	// has already been replaced by Kitty image protocol strings.
    	// We do not want parser.MathJax here as it would try to process those protocol strings.
    	p := parser.New() // No extensions, especially not MathJax
    	docNode := p.Parse([]byte(line))

    	var sb strings.Builder
    	renderMarkdownAST(docNode, &sb)
    	
	// Get the raw string from the builder
	// Apply Markdown formatting to the final processed string.
	result := sb.String()
	
	// Process any escape characters that might appear in the output
	result = strings.ReplaceAll(result, "\\n", "\n") // Replace escaped newlines
	result = strings.ReplaceAll(result, "\\%", "%")  // Replace escaped percent signs
	
	// Clean up any unexpected trailing characters
	result = strings.TrimSuffix(result, "%")  
	result = strings.TrimSuffix(result, "\x00")
	
	return result
    }

    func main() {
	// Command-line flags
	colourFlag := flag.String("colour", "white", "Set LaTeX text colour (e.g., red, #00FF00).")
	cFlag := flag.String("c", "", "Short alias for --colour. Overrides --colour if set.")
	sizeFlag := flag.Int("size", 0, "Target terminal rows for LaTeX images (0 for default: 1 for inline, auto for display).")
	sFlag := flag.Int("s", 0, "Short alias for --size.")
	dpiFlag := flag.Int("dpi", 300, "Set DPI for rendering LaTeX images.")
	dFlag := flag.Int("d", 0, "Short alias for --dpi. Overrides --dpi if set (and not 0).")
	renderAllLatexFlag := flag.Bool("render-all-latex", false, "Render entire input as a single LaTeX document/image.")
	lFlag := flag.Bool("l", false, "Short alias for --render-all-latex.")
	debugFlag := flag.Bool("debug", false, "Enable verbose debug output.")
	dDebugFlag := flag.Bool("D", false, "Short alias for --debug.")

	flag.Parse() // Parse all flags first

	// Set default color to white and apply overrides if specified
	effectiveColor := "white"
	if *cFlag != "" {
		effectiveColor = *cFlag
	} else if *colourFlag != "white" {
		effectiveColor = *colourFlag
	}

	// Determine if short flags -s and -d were explicitly set
	var sFlagSet, dFlagSet bool
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "s":
			sFlagSet = true
		case "d":
			dFlagSet = true
		}
	})

	// Determine the effective size to use
	effectiveSize := *sizeFlag
	if sFlagSet { // If -s was explicitly provided on the command line, it takes precedence
		effectiveSize = *sFlag
	}
	if effectiveSize < 0 { // Treat negative size as default (0)
		effectiveSize = 0
	}

	// Determine the effective DPI to use
	effectiveDPI := *dpiFlag
	if dFlagSet { // If -d was explicitly provided on the command line, it takes precedence
		effectiveDPI = *dFlag
	}
	if effectiveDPI <= 0 { // Ensure DPI is positive, default to 300 if invalid or explicitly set to 0 or less
		effectiveDPI = 300
	}

	isRenderAllLatexMode := *renderAllLatexFlag || *lFlag
	isDebugMode := *debugFlag || *dDebugFlag

	if isDebugMode {
		fmt.Fprintln(os.Stderr, "DEBUG: dml starting")
		fmt.Fprintln(os.Stderr, "DEBUG: Flags parsed.")
		fmt.Fprintf(os.Stderr, "DEBUG: isRenderAllLatexMode: %v\n", isRenderAllLatexMode)
	}

	if isRenderAllLatexMode {
		if isDebugMode {
			if isDebugMode {
				fmt.Fprintln(os.Stderr, "DEBUG: Reading standard input (full document) for render-all-latex mode...")
			}
		}
		// Read all of stdin into a single string for full LaTeX mode
		inputBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
			os.Exit(1)
		}
		inputString := string(inputBytes)
		if isDebugMode {
			if isDebugMode {
				fmt.Fprintf(os.Stderr, "DEBUG: Finished reading input (%d bytes).\n", len(inputBytes))
			}
		}

		// Full LaTeX rendering mode
		// Preprocess \[...\] and \(...\) to $$...$$ and $...$ for correct math parsing
		preprocessed := displayMathBracket.ReplaceAllStringFunc(inputString, func(match string) string {
			content := strings.TrimSpace(match[2 : len(match)-2])
			return "$" + content + "$"
		})
		preprocessed = inlineMathParen.ReplaceAllStringFunc(preprocessed, func(match string) string {
			content := strings.TrimSpace(match[2 : len(match)-2])
			return "$" + content + "$"
		})

		// Enable MathJax and other common extensions for parsing
		// The parser.MathJax extension is key here as it will create ast.Math and ast.MathBlock nodes.
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.MathJax)
		docNode := p.Parse([]byte(preprocessed))

		var latexBodyBuilder strings.Builder
		generateLatexFromAST(docNode, &latexBodyBuilder)
		latexBody := latexBodyBuilder.String()

		img, renderErr := renderFullLatexDocument(latexBody, effectiveColor, effectiveDPI)
		if renderErr != nil {
			fmt.Fprintf(os.Stderr, "Error in full LaTeX rendering mode: %v\\n", renderErr)
			// In full render mode, if LaTeX fails, print the original input so user can debug
			fmt.Print(inputString)
			os.Exit(1)
		}

		kittyStr, kittyErr := kittyInline(img, true, effectiveSize)
		if kittyErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for full document: %v\\n", kittyErr)
			// If Kitty protocol generation fails, print original input
			fmt.Print(inputString)
			os.Exit(1)
		}
		fmt.Print(kittyStr)

	} else {
		if isDebugMode {
			if isDebugMode {
				fmt.Fprintln(os.Stderr, "DEBUG: Entering standard processing mode (line-by-line streaming with state).")
			}
		}
		// Streaming processing mode with state for multi-line math
		// Reads input line by line and processes incrementally.
		// Uses a state machine to buffer multi-line display math blocks.

		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(os.Stdout) // Use a buffered writer for output flushing

		var mathBuffer strings.Builder // Buffer for collecting multi-line math content
		inDisplayMath := false         // State flag

		// Set environment variable for renderMath debug flag if we're in debug mode
		if isDebugMode {
			os.Setenv("DML_DEBUG", "1")
			fmt.Fprintln(os.Stderr, "DEBUG: Starting line-by-line input reading and processing loop...")
		}

		// Regexes to find delimiters specifically for streaming processing.
		// These look for the *start* and *end* of delimiters on the current line.
		// They are different from the global regexes which assume the whole input is available.
		startDisplayMathRegex := regexp.MustCompile(`(?:^|\s)\$\$|\\\[`)
		endDisplayMathRegex := regexp.MustCompile(`\$\$(?:$|\s)|\\\]`)

		// Set environment variable for renderMath debug flag if we're in debug mode
		if isDebugMode {
		    os.Setenv("DML_DEBUG", "1")
			fmt.Fprintln(os.Stderr, "DEBUG: Using regexes for display math detection")
		}
		// Note: This simplified approach won't handle inline math inside display math,
		// or multiple display math blocks on one line perfectly if delimiters are mixed.
		// It prioritizes streaming output for simple cases.

		for {
			inputLine, err := reader.ReadString('\n')

			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading input line: %v\n", err)
				writer.Flush() // Flush any pending output
				os.Exit(1)
			}

			// Determine if this is the last line
			isLastLine := (err == io.EOF)

			if inDisplayMath {
				// We are inside a display math block
				if isDebugMode {
					fmt.Fprintf(os.Stderr, "DEBUG: In normal mode, processing line: %s\n", strings.TrimSpace(inputLine))
				}

				endMatchIdx := endDisplayMathRegex.FindStringIndex(inputLine)

				if endMatchIdx != nil {
					// Found the closing delimiter on this line
					if isDebugMode {
						fmt.Fprintln(os.Stderr, "DEBUG: Found display math end delimiter.")
					}
					mathBuffer.WriteString(inputLine[:endMatchIdx[0]]) // Add content before the delimiter

					mathContent := mathBuffer.String()
					mathBuffer.Reset() // Clear the buffer
					inDisplayMath = false // Exit display math state

					// Render the collected math content
					if isDebugMode {
						fmt.Fprintf(os.Stderr, "DEBUG: Attempting to render display math (length: %d chars)\n", len(mathContent))
					}
					img, renderErr := renderMath(mathContent, effectiveColor, true, effectiveDPI)
					if renderErr != nil {
						if isDebugMode {
							fmt.Fprintf(os.Stderr, "DEBUG: Attempting to render display math (length: %d chars)\n", len(mathContent))
						}
						fmt.Fprintf(os.Stderr, "ERROR: Rendering display math failed: %v\n", renderErr)
						// On error, print the un-rendered content as text
						writer.WriteString("$$")
						writer.WriteString(mathContent)
						writer.WriteString("$$\n") // Add newline if it's missing from content
					} else {
						if isDebugMode {
							fmt.Fprintln(os.Stderr, "DEBUG: Math rendering successful, generating Kitty protocol")
						}
						kittyStr, kittyErr := kittyInline(img, true, effectiveSize)
						if kittyErr != nil {
							fmt.Fprintf(os.Stderr, "ERROR: Generating Kitty protocol failed: %v\n", kittyErr)
							// On error, print the un-rendered content as text
							writer.WriteString("$$")
							writer.WriteString(mathContent)
							writer.WriteString("$$\n") // Add newline if it's missing from content
						} else {
							if isDebugMode {
								fmt.Fprintln(os.Stderr, "DEBUG: Successfully generated Kitty protocol for display math")
							}
							writer.WriteString(kittyStr) // Write the rendered image protocol
						}
					}

					// Process the rest of the line after the closing delimiter
					remainingLine := inputLine[endMatchIdx[1]:]
					if len(remainingLine) > 0 {
						// Process remaining part of the line as normal text (inline math, markdown)
						if isDebugMode {
							fmt.Fprintf(os.Stderr, "DEBUG: Processing remaining line after display math: %s\n", strings.TrimSpace(remainingLine))
						}
						processedRemaining := inlineMath.ReplaceAllStringFunc(remainingLine, func(match string) string {
							content := strings.TrimSpace(match[1 : len(match)-1])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
								return match
							}
							return kStr
						})
						processedRemaining = inlineMathParen.ReplaceAllStringFunc(processedRemaining, func(match string) string {
							content := strings.TrimSpace(match[2 : len(match)-2])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
								return match
							}
							return kStr
						})
						finalRemainingOutput := applyMarkdownFormatting(processedRemaining)
						writer.WriteString(finalRemainingOutput)
					}

				} else {
					// No closing delimiter yet, just buffer the line
					mathBuffer.WriteString(inputLine)
					if isDebugMode {
						fmt.Fprintln(os.Stderr, "DEBUG: Appended line to math buffer.")
					}
				}

			} else {
				// We are in normal text mode
				if isDebugMode {
					fmt.Fprintf(os.Stderr, "DEBUG: In normal mode, processing line: %s\n", strings.TrimSpace(inputLine))
				}

				startMatchIdx := startDisplayMathRegex.FindStringIndex(inputLine)
				endMatchIdx := endDisplayMathRegex.FindStringIndex(inputLine) // Check for same-line closing

				if startMatchIdx != nil && (endMatchIdx == nil || endMatchIdx[0] < startMatchIdx[0]) {
					// Found starting delimiter for a multi-line block (and no closing before it)
					if isDebugMode {
						fmt.Fprintln(os.Stderr, "DEBUG: Found display math start delimiter. Switching to math state.")
					}
					// Process content *before* the delimiter as normal text
					beforeDelimiter := inputLine[:startMatchIdx[0]]
					if len(beforeDelimiter) > 0 {
						if isDebugMode {
							fmt.Fprintf(os.Stderr, "DEBUG: Processing text before delimiter: %s\n", strings.TrimSpace(beforeDelimiter))
						}
						processedBefore := inlineMath.ReplaceAllStringFunc(beforeDelimiter, func(match string) string {
							content := strings.TrimSpace(match[1 : len(match)-1])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
								return match
							}
							return kStr
						})
						processedBefore = inlineMathParen.ReplaceAllStringFunc(processedBefore, func(match string) string {
							content := strings.TrimSpace(match[2 : len(match)-2])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
								return match
							}
							return kStr
						})

						finalBeforeOutput := applyMarkdownFormatting(processedBefore)
						writer.WriteString(finalBeforeOutput)
					}

					// Start buffering from the content *after* the delimiter on this line
					mathBuffer.WriteString(inputLine[startMatchIdx[1]:])
					inDisplayMath = true // Enter display math state
					if isDebugMode {
						fmt.Fprintln(os.Stderr, "DEBUG: Started buffering math content.")
					}

				} else if startMatchIdx != nil && endMatchIdx != nil && startMatchIdx[0] < endMatchIdx[0] {
					// Found both start and end delimiters on the same line (single-line display math)
					if isDebugMode {
						fmt.Fprintln(os.Stderr, "DEBUG: Found single-line display math.")
					}
					// Process content *before* the start delimiter
					beforeDelimiter := inputLine[:startMatchIdx[0]]
					if len(beforeDelimiter) > 0 {
						if isDebugMode {
							fmt.Fprintf(os.Stderr, "DEBUG: Processing text before single-line math: %s\n", strings.TrimSpace(beforeDelimiter))
						}
						processedBefore := inlineMath.ReplaceAllStringFunc(beforeDelimiter, func(match string) string {
							content := strings.TrimSpace(match[1 : len(match)-1])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
								return match
							}
							return kStr
						})
						processedBefore = inlineMathParen.ReplaceAllStringFunc(processedBefore, func(match string) string {
							content := strings.TrimSpace(match[2 : len(match)-2])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
								return match
							}
							return kStr
						})

						finalBeforeOutput := applyMarkdownFormatting(processedBefore)
						writer.WriteString(finalBeforeOutput)
					}

					// Extract and process the math content
					mathContent := inputLine[startMatchIdx[1]:endMatchIdx[0]]
					if isDebugMode {
						fmt.Fprintf(os.Stderr, "DEBUG: Rendering single-line display math: %s\n", strings.TrimSpace(mathContent))
					}
					img, renderErr := renderMath(mathContent, effectiveColor, true, effectiveDPI)
					if renderErr != nil {
						fmt.Fprintf(os.Stderr, "ERROR: Rendering display math failed: %v\n", renderErr)
						// On error, print the un-rendered content as text
						writer.WriteString("$$")
						writer.WriteString(mathContent)
						writer.WriteString("$$\n") // Add newline if it's missing from content
					} else {
						kittyStr, kittyErr := kittyInline(img, true, effectiveSize)
						if kittyErr != nil {
							fmt.Fprintf(os.Stderr, "ERROR: Generating Kitty protocol failed: %v\n", kittyErr)
							// On error, print the un-rendered content as text
							writer.WriteString("$$")
							writer.WriteString(mathContent)
							writer.WriteString("$$\n") // Add newline if it's missing from content
						} else {
							writer.WriteString(kittyStr) // Write the rendered image protocol
						}
					}

					// Process content *after* the end delimiter
					afterDelimiter := inputLine[endMatchIdx[1]:]
					if len(afterDelimiter) > 0 {
						if isDebugMode {
							fmt.Fprintf(os.Stderr, "DEBUG: Processing text after single-line math: %s\n", strings.TrimSpace(afterDelimiter))
						}
						processedAfter := inlineMath.ReplaceAllStringFunc(afterDelimiter, func(match string) string {
							content := strings.TrimSpace(match[1 : len(match)-1])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\\n", content, kErr)
								return match
							}
							return kStr
						})
						processedAfter = inlineMathParen.ReplaceAllStringFunc(processedAfter, func(match string) string {
							content := strings.TrimSpace(match[2 : len(match)-2])
							if content == "" { return match }
							img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
							if rErr != nil {
								fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\\n", content, rErr)
								return match
							}
							kStr, kErr := kittyInline(img, false, effectiveSize)
							if kErr != nil {
								fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\\n", content, kErr)
								return match
							}
							return kStr
						})

						finalAfterOutput := applyMarkdownFormatting(processedAfter)
						writer.WriteString(finalAfterOutput)
					}

				} else {
					// No display math delimiters found on this line.
					// Process for inline math and markdown as before.
					processedLine := inlineMath.ReplaceAllStringFunc(inputLine, func(match string) string {
						content := strings.TrimSpace(match[1 : len(match)-1])
						if content == "" { return match }
						img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
						if rErr != nil {
							fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\\n", content, rErr)
							return match
						}
						kStr, kErr := kittyInline(img, false, effectiveSize)
						if kErr != nil {
							fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\\n", content, kErr)
							return match
						}
						return kStr
					})
					processedLine = inlineMathParen.ReplaceAllStringFunc(processedLine, func(match string) string {
						content := strings.TrimSpace(match[2 : len(match)-2])
						if content == "" { return match }
						img, rErr := renderMath(content, effectiveColor, false, effectiveDPI)
						if rErr != nil {
							fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\\n", content, rErr)
							return match
						}
						kStr, kErr := kittyInline(img, false, effectiveSize)
						if kErr != nil {
							fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\\n", content, kErr)
							return match
						}
						return kStr
					})

					// Apply Markdown formatting to the processed line.
					finalLineOutput := applyMarkdownFormatting(processedLine)
					
					// Remove any trailing special characters that might appear
					finalLineOutput = strings.TrimSuffix(finalLineOutput, "%")
					finalLineOutput = strings.TrimSuffix(finalLineOutput, "\x00")
					
					// Make sure we keep newlines as is
					if strings.HasSuffix(processedLine, "\n") && !strings.HasSuffix(finalLineOutput, "\n") {
						finalLineOutput += "\n"
					}

					// Write the processed line to the buffered writer.
					writer.WriteString(finalLineOutput)
				}
			}

			// Flush the writer immediately after processing a line (or block end).
			writer.Flush()

			// If the error was EOF, it means we just processed the last line.
			// The loop condition should break after this iteration.
			if isLastLine {
				if isDebugMode {
					fmt.Fprintln(os.Stderr, "DEBUG: Finished processing line before EOF. Exiting loop after flush.")
				}
				break
			}
		}

		// Handle case where input ended while still inside a display math block
		if inDisplayMath {
			if isDebugMode {
				fmt.Fprintln(os.Stderr, "DEBUG: Warning: Reached EOF while still inside a display math block. Outputting buffered content as plain text.")
			}
			writer.WriteString("$$") // Output the start delimiter that wasn't closed
			writer.WriteString(mathBuffer.String()) // Output the buffered content
			// No closing delimiter to output
			writer.Flush()
		}

		if isDebugMode {
			fmt.Fprintln(os.Stderr, "DEBUG: Output streaming finished.")
		}
		}
	
	// Final debug messages if debug mode is enabled
	if isDebugMode {
		fmt.Fprintf(os.Stderr, "DEBUG: dml execution completed. If math rendering issues occurred, check for LaTeX or convert errors.")
		fmt.Fprintln(os.Stderr, "DEBUG: dml exiting.")
	}
}
