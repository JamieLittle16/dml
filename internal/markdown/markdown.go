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
	// Default handling for other common nodes: just process their children.
	case *ast.List, *ast.ListItem, *ast.Link, *ast.Image,
		*ast.BlockQuote, *ast.HorizontalRule, *ast.HTMLBlock, *ast.HTMLSpan,
		*ast.Table, *ast.TableCell, *ast.TableHeader, *ast.TableRow:
		processChildren(n) // For now, just try to render their content
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
	// For common containers, just recurse on children
	case *ast.Document, *ast.Paragraph, *ast.List, *ast.ListItem, *ast.Link, *ast.Image,
		*ast.Code, *ast.CodeBlock, *ast.BlockQuote, *ast.Heading,
		*ast.HorizontalRule, *ast.HTMLBlock, *ast.HTMLSpan,
		*ast.Table, *ast.TableCell, *ast.TableHeader, *ast.TableRow,
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
	// Standard parser, no special extensions
	p := parser.New() // No extensions, especially not MathJax
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