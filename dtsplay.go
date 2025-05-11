package main

import (
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
	"github.com/gomarkdown/markdown/parser" // Renamed alias to mdparser to avoid conflict
)

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
	// Default handling for other common nodes: just process their children.
	// This will extract text, bold, italic, and math from within them.
	// More specific LaTeX conversions could be added here for lists, headings, etc.
	case *ast.List, *ast.ListItem, *ast.Link, *ast.Image,
		*ast.Code, *ast.CodeBlock, *ast.BlockQuote, *ast.Heading,
		*ast.HorizontalRule, *ast.HTMLBlock, *ast.HTMLSpan,
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
	var colorCmd string
	if strings.HasPrefix(color, "#") {
		colorCmd = fmt.Sprintf(`\color[HTML]{%s}`, strings.TrimPrefix(color, "#"))
	} else {
		colorCmd = fmt.Sprintf(`\color{%s}`, color)
	}

	// Using standalone with varwidth=\maxdimen to try and fit content.
	// Added amssymb, amsfonts, inputenc, fontenc, lmodern for better general text and math.
	texTemplate := `\documentclass[preview,varwidth=\maxdimen]{standalone}
\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{amsfonts}
\usepackage{xcolor}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{lmodern}
% Any other packages needed for Markdown features would go here
\begin{document}
%s %s
\end{document}`
	// Note: The color command is now applied once at the beginning of the document body.
	// The body itself might contain $$...$$ or $...$ which will inherit this color.
	tex := fmt.Sprintf(texTemplate, colorCmd, latexBody)

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
	cmd = exec.Command("convert", "-density", fmt.Sprintf("%d", dpi), pdfFile, pngFile)
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
    	var colorCmd string
    	if strings.HasPrefix(color, "#") {
    		colorCmd = fmt.Sprintf(`\color[HTML]{%s}`, strings.TrimPrefix(color, "#"))
    	} else {
    		colorCmd = fmt.Sprintf(`\color{%s}`, color)
    	}

    	var texTemplate string
    	if isDisplay {
    		texTemplate = `\documentclass[preview]{standalone}
	\usepackage{amsmath}
	\usepackage{xcolor}
	\begin{document}
	%s\[%s\]
	\end{document}`
    	} else {
    		texTemplate = `\documentclass[preview]{standalone}
	\usepackage{amsmath}
	\usepackage{xcolor}
	\begin{document}
	%s$%s$
	\end{document}`
    	}
    	tex := fmt.Sprintf(texTemplate, colorCmd, latex)

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
    cmd = exec.Command("convert", "-density", fmt.Sprintf("%d", dpi), pdfFile, pngFile)
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

// kittyInline generates the Kitty graphics protocol string for the given image bytes.
func kittyInline(img []byte, isDisplayMath bool, userTargetRows int) (string, error) {
    var sb strings.Builder
    opts := rasterm.KittyImgOpts{
        // Available fields in v1.1.1:
        // SrcX, SrcY, SrcWidth, SrcHeight, CellOffsetX, CellOffsetY,
        // DstCols, DstRows, ZIndex, ImageId, ImageNo, PlacementId
        // TransferMode and UseWindowSize are not available in this version.
    }

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
        return "", err
    }
    kittyStr := sb.String()
    if isDisplayMath {
        kittyStr += "\n" // Add a newline after display math images
    }
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
    	return sb.String()
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
	flag.Parse()

	// Determine the effective color to use
	effectiveColor := *colourFlag
	if *cFlag != "" {
		effectiveColor = *cFlag
	}

	// Determine the effective size to use
	effectiveSize := *sizeFlag
	if *sFlag != 0 { // If -s was used (and not explicitly set to its default 0)
		effectiveSize = *sFlag
	}
	if effectiveSize < 0 { // Treat negative size as default
		effectiveSize = 0
	}

	// Determine the effective DPI to use
	effectiveDPI := *dpiFlag
	if *dFlag != 0 { // If -d was used and not its default 0
		effectiveDPI = *dFlag
	}
	if effectiveDPI <= 0 { // Ensure DPI is positive, default to 300 if invalid
		effectiveDPI = 300
	}

	isRenderAllLatexMode := *renderAllLatexFlag || *lFlag

	// Read all of stdin into a single string
	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
		os.Exit(1)
	}
	inputString := string(inputBytes)

	if isRenderAllLatexMode {
		// Full LaTeX rendering mode
		// Enable MathJax and other common extensions for parsing
		// The parser.MathJax extension is key here as it will create ast.Math and ast.MathBlock nodes.
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.MathJax)
		docNode := p.Parse(inputBytes)

		var latexBodyBuilder strings.Builder
		generateLatexFromAST(docNode, &latexBodyBuilder)
		latexBody := latexBodyBuilder.String()

		// For debugging the generated LaTeX:
		// fmt.Fprintf(os.Stderr, "--- Generated LaTeX Body ---\n%s\n--------------------------\n", latexBody)

		img, renderErr := renderFullLatexDocument(latexBody, effectiveColor, effectiveDPI)
		if renderErr != nil {
			fmt.Fprintf(os.Stderr, "Error in full LaTeX rendering mode: %v\n", renderErr)
			// Fallback: print original input or a more specific error message?
			// For now, print original input on error to avoid losing it.
			fmt.Print(inputString)
			os.Exit(1)
		}

		// Treat the entire document image as "display math" for sizing and newline.
		// Use effectiveSize for target rows.
		kittyStr, kittyErr := kittyInline(img, true, effectiveSize)
		if kittyErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for full document: %v\n", kittyErr)
			fmt.Print(inputString) // Fallback
			os.Exit(1)
		}
		fmt.Print(kittyStr) // Use Print not Println as kittyInline adds \n for display mode

	} else {
		// Original processing mode (per-expression LaTeX, Markdown for bold/italic)
		line := inputString
		// Step 1: Process LaTeX for rendering to Kitty protocol strings
		// Process display-math blocks ($$...$$) first
		line = displayMath.ReplaceAllStringFunc(line, func(match string) string {
			content := match[2 : len(match)-2] // Extract LaTeX content for display math
			img, err := renderMath(content, effectiveColor, true, effectiveDPI)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering display math ('%s'): %v\n", content, err)
				return match
			}
			kittyStr, err := kittyInline(img, true, effectiveSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for display math ('%s'): %v\n", content, err)
				return match // Return original LaTeX string if Kitty protocol generation fails
			}
			return kittyStr // Replace LaTeX with Kitty image protocol string
		})

		// Process display-math blocks (\[...\])
		line = displayMathBracket.ReplaceAllStringFunc(line, func(match string) string {
			// Extract content from between \[ and \]
			content := strings.TrimSpace(match[2 : len(match)-2])
			if content == "" { return match }
			img, err := renderMath(content, effectiveColor, true, effectiveDPI) // Corrected: isDisplay is true
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering display math ('%s'): %v\n", content, err) // Corrected: error message
				return match
			}
			kittyStr, err := kittyInline(img, true, effectiveSize) // Corrected: isDisplayMath is true
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for display math ('%s'): %v\n", content, err) // Corrected: error message
				return match
			}
			return kittyStr
		})

		// Process inline-math snippets ($...$)
		line = inlineMath.ReplaceAllStringFunc(line, func(match string) string {
			// Extract content from between $ and $
			content := strings.TrimSpace(match[1 : len(match)-1])
			if content == "" { return match }
			img, err := renderMath(content, effectiveColor, false, effectiveDPI)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, err)
				return match
			}
			kittyStr, err := kittyInline(img, false, effectiveSize) // Corrected: add effectiveSize
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, err)
				return match
			}
			return kittyStr
		})

		// Process inline-math snippets (\\(...\\))
		line = inlineMathParen.ReplaceAllStringFunc(line, func(match string) string {
			// Extract content from between \( and \)
			content := strings.TrimSpace(match[2 : len(match)-2])
			if content == "" { return match }
			img, err := renderMath(content, effectiveColor, false, effectiveDPI)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, err) // Corrected: error message
				return match
			}
			kittyStr, err := kittyInline(img, false, effectiveSize) // Corrected: isDisplayMath is false
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, err) // Corrected: error message
				return match
			}
			return kittyStr
		})

		// Step 2: Apply Markdown formatting (italics, bold) to the line,
		// which now contains text and Kitty protocol strings from rendered LaTeX.
		// Note: The `parser.MathJax` extension is NOT used for applyMarkdownFormatting
		// because LaTeX has already been replaced by Kitty sequences.
		line = applyMarkdownFormatting(line)

		// Print the fully processed line.
		fmt.Println(line)
	}
	// Note: scanner.Err() is not applicable here as we read all input at once.
}
