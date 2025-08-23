// Package markdown provides markdown processing functionality for DML
package markdown

import (
	"strings"

	"dml/internal/latex"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// GenerateLatexFromAST traverses a markdown AST and generates LaTeX output
func GenerateLatexFromAST(node ast.Node, sb *strings.Builder) {
	if node == nil {
		return
	}

	// Helper to process children
	processChildren := func(n ast.Node) {
		for _, child := range n.GetChildren() {
			GenerateLatexFromAST(child, sb)
		}
	}

	switch n := node.(type) {
	case *ast.Document:
		processChildren(n)
	case *ast.Text:
		sb.WriteString(latex.EscapeLaTeX(string(n.Literal)))
	case *ast.Emph: // *italic* or _italic_
		sb.WriteString(`\textit{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Strong: // **bold** or __bold__
		sb.WriteString(`\textbf{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Del: // ~~strikethrough~~
		sb.WriteString(`\sout{`) // Requires ulem package
		processChildren(n)
		sb.WriteString(`}`)
	// TODO: Add cases for other GFM extensions:
	// - Task list items (checkbox symbols)
	// - Emoji rendering (emoji package or unicode)
	// - Autolinks (href package)
	// - Definition lists (description environment)
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
		sb.WriteString(latex.EscapeLaTeX(string(n.Literal)))
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
	// Enhanced Markdown features - scaffolding for future implementation
	case *ast.List:
		// TODO: Implement proper LaTeX list rendering (itemize/enumerate)
		// For now, just process children with basic spacing
		sb.WriteString("\n") // Add spacing before list
		processChildren(n)
		sb.WriteString("\n") // Add spacing after list
	case *ast.ListItem:
		// TODO: Add proper list item bullets/numbers based on parent list type
		sb.WriteString("• ") // Basic bullet for now
		processChildren(n)
		sb.WriteString("\\\\\n") // Line break after item
	case *ast.Link:
		// TODO: Implement proper LaTeX link handling with href package
		// For now, render as: text (url)
		processChildren(n) // Render link text
		if n.Destination != nil {
			sb.WriteString(" (")
			sb.WriteString(latex.EscapeLaTeX(string(n.Destination)))
			sb.WriteString(")")
		}
	case *ast.Image:
		// TODO: Implement image handling - for now render alt text
		if n.Title != nil {
			sb.WriteString("[Image: ")
			sb.WriteString(latex.EscapeLaTeX(string(n.Title)))
			sb.WriteString("]")
		} else {
			sb.WriteString("[Image]")
		}
		// Note: n.Destination contains the image URL for future implementation
	case *ast.BlockQuote:
		// TODO: Implement proper LaTeX quote environment
		sb.WriteString("\\begin{quote}\n")
		processChildren(n)
		sb.WriteString("\\end{quote}\n")
	case *ast.Table:
		// TODO: Implement full LaTeX table with tabular environment
		sb.WriteString("\\begin{center}\n\\begin{tabular}{|")
		// For now, just create a basic table structure
		sb.WriteString("l|l|l|") // Placeholder for 3 columns
		sb.WriteString("}\n\\hline\n")
		processChildren(n)
		sb.WriteString("\\hline\n\\end{tabular}\n\\end{center}\n")
	case *ast.TableHeader:
		// TODO: Implement table header formatting
		processChildren(n)
		sb.WriteString(" \\\\ \\hline\n")
	case *ast.TableRow:
		// TODO: Implement proper table row formatting
		processChildren(n)
		sb.WriteString(" \\\\\n")
	case *ast.TableCell:
		// TODO: Implement proper cell content and alignment
		processChildren(n)
		sb.WriteString(" & ") // Column separator
	case *ast.HorizontalRule:
		sb.WriteString("\\hrule\n")
	case *ast.HTMLBlock, *ast.HTMLSpan:
		// TODO: Implement HTML passthrough or conversion
		processChildren(n) // For now, just render content
	default:
		// For any other unhandled node type, if it's a container, process its children.
		children := n.GetChildren()
		for _, child := range children {
			GenerateLatexFromAST(child, sb)
		}
	}
}

// RenderMarkdownAST recursively traverses the AST and builds a string with ANSI codes
func RenderMarkdownAST(node ast.Node, sb *strings.Builder) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.Text:
		sb.Write(n.Literal)
	case *ast.Emph: // Handles *italic* and _italic_
		sb.WriteString("\x1b[3m") // ANSI Italic on
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[23m") // ANSI Italic off
	case *ast.Strong: // Handles **bold** and __bold__
		sb.WriteString("\x1b[1m") // ANSI Bold on
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[22m") // ANSI Bold off
	case *ast.Code: // Inline code
		// TODO: Add syntax highlighting support for code spans
		sb.WriteString("\x1b[7m") // ANSI Reverse video for code
		sb.Write(n.Literal)
		sb.WriteString("\x1b[27m") // ANSI Reverse video off
	case *ast.CodeBlock: // Code blocks
		// TODO: Add syntax highlighting based on language info
		sb.WriteString("\x1b[7m") // ANSI Reverse video for code
		sb.WriteString("\n")
		sb.Write(n.Literal)
		sb.WriteString("\x1b[27m") // ANSI Reverse video off
		sb.WriteString("\n")
	case *ast.Heading:
		// TODO: Add different styling for different heading levels
		sb.WriteString("\x1b[1m\x1b[4m") // Bold + Underline for headings
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[22m\x1b[24m") // Bold off + Underline off
		sb.WriteString("\n")
	// Enhanced Markdown features - scaffolding for terminal rendering
	case *ast.List:
		// TODO: Implement proper list formatting with indentation
		sb.WriteString("\n") // Add spacing before list
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\n") // Add spacing after list
	case *ast.ListItem:
		// TODO: Add proper bullet/number based on list type and nesting level
		if n.ListFlags&ast.ListTypeOrdered != 0 {
			sb.WriteString("1. ") // Placeholder for ordered lists
		} else {
			sb.WriteString("• ") // Bullet for unordered lists
		}
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\n")
	case *ast.Link:
		// TODO: Add clickable terminal links where supported (OSC 8)
		sb.WriteString("\x1b[4m\x1b[34m") // Underline + Blue for links
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[24m\x1b[39m") // Underline off + Default color
		// TODO: Add URL display or click handling
	case *ast.Image:
		// TODO: Implement image placeholder or inline image display
		sb.WriteString("\x1b[33m[Image") // Yellow text for image placeholder
		if len(n.Children) > 0 {
			sb.WriteString(": ")
			for _, child := range n.GetChildren() {
				RenderMarkdownAST(child, sb)
			}
		}
		sb.WriteString("]\x1b[39m") // Close bracket + default color
	case *ast.BlockQuote:
		// TODO: Implement proper blockquote formatting with bars
		sb.WriteString("\x1b[36m") // Cyan color for quotes
		sb.WriteString("│ ")       // Quote bar
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[39m") // Default color
		sb.WriteString("\n")
	case *ast.Table:
		// TODO: Implement full table rendering with borders and alignment
		sb.WriteString("\x1b[1m") // Bold for table
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[22m") // Bold off
	case *ast.Del: // ~~strikethrough~~
		sb.WriteString("\x1b[9m") // ANSI Strikethrough on
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[29m") // ANSI Strikethrough off
	// TODO: Add cases for other GFM extensions:
	// - Task list items (checked/unchecked checkboxes)
	// - Emoji rendering (:emoji: syntax)
	// - Autolinks
	// - Definition lists
		sb.WriteString("\n")
	case *ast.TableHeader:
		// TODO: Implement table header with proper formatting
		sb.WriteString("\x1b[4m") // Underline for headers
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[24m") // Underline off
		sb.WriteString("\n")
	case *ast.TableRow:
		// TODO: Implement table row formatting
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\n")
	case *ast.TableCell:
		// TODO: Implement proper cell formatting with padding and alignment
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString(" | ") // Column separator
	case *ast.HorizontalRule:
		// TODO: Add customizable horizontal rule rendering
		sb.WriteString("\x1b[2m") // Dim
		sb.WriteString("─────────────────────────────────────")
		sb.WriteString("\x1b[22m") // Dim off
		sb.WriteString("\n")
	// For common containers, just recurse on children
	case *ast.Document, *ast.Paragraph,
		*ast.HTMLBlock, *ast.HTMLSpan,
		*ast.Math, *ast.MathBlock: // Include Math nodes in case parser produces them
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
	default:
		// For any other unhandled node type, process its children
		children := n.GetChildren()
		for _, child := range children {
			RenderMarkdownAST(child, sb)
		}
	}
}

// ApplyFormatting parses a line of Markdown text and converts basic styling to ANSI codes
func ApplyFormatting(line string) string {
	// Create parser with GFM extensions for enhanced markdown support
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	// Enable specific GFM extensions
	extensions |= parser.Tables | parser.Strikethrough | parser.FencedCode | parser.Autolink
	
	p := parser.NewWithExtensions(extensions)
	docNode := p.Parse([]byte(line))

	var sb strings.Builder
	RenderMarkdownAST(docNode, &sb)

	// Get the raw string from the builder
	result := sb.String()

	// Process any escape characters that might appear in the output
	result = strings.ReplaceAll(result, "\\n", "\n") // Replace escaped newlines
	result = strings.ReplaceAll(result, "\\%", "%")  // Replace escaped percent signs

	// Clean up any unexpected trailing characters
	result = strings.TrimSuffix(result, "%")
	result = strings.TrimSuffix(result, "\x00")

	return result
}