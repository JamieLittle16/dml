// Package markdown provides markdown processing functionality for DML
package markdown

import (
	"os"
	"strconv"
	"strings"

	"dml/internal/latex"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// RenderState holds state for recursive rendering (list indices, prefixes, etc.)
type RenderState struct {
	ListDepth      int
	OrderedIndex   int
	InBlockquote   bool
	BlockquotePrefix string
}

// Helper functions for ANSI rendering

// getTerminalWidth tries to get terminal width, falls back to 80 chars
func getTerminalWidth() int {
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}
	return 80 // Default fallback
}

// renderHorizontalRule creates a line of characters for terminal display
func renderHorizontalRule() string {
	width := getTerminalWidth()
	if width > 40 {
		width = 40 // Cap at reasonable length
	}
	return strings.Repeat("─", width) // Unicode box drawing character, fallback to hyphens if needed
}

// renderListPrefix returns the appropriate prefix for list items
func renderListPrefix(isList *ast.List, itemIndex int) string {
	if isList.ListFlags&ast.ListTypeOrdered != 0 {
		return strconv.Itoa(itemIndex+1) + ". "
	}
	// Try Unicode bullet, fallback to ASCII
	return "• "
}

// renderBlockquotePrefix returns the prefix for blockquote lines
func renderBlockquotePrefix() string {
	// Try Unicode box drawing character, fallback to ASCII
	return "│ "
}

// SimpleTable represents a basic table structure for ASCII rendering
type SimpleTable struct {
	Headers []string
	Rows    [][]string
	ColWidths []int
}

// collectTableData extracts table data from AST nodes
func collectTableData(tableNode *ast.Table) *SimpleTable {
	table := &SimpleTable{}
	
	for _, child := range tableNode.GetChildren() {
		switch node := child.(type) {
		case *ast.TableHeader:
			// Collect headers from the first (and usually only) row in the header
			for _, headerChild := range node.GetChildren() {
				if rowNode, ok := headerChild.(*ast.TableRow); ok {
					for _, cell := range rowNode.GetChildren() {
						if cellNode, ok := cell.(*ast.TableCell); ok {
							var cellText strings.Builder
							collectTextFromNode(cellNode, &cellText)
							table.Headers = append(table.Headers, strings.TrimSpace(cellText.String()))
						}
					}
					break // Only process the first row in header
				}
			}
		case *ast.TableBody:
			// Collect body rows
			for _, row := range node.GetChildren() {
				if rowNode, ok := row.(*ast.TableRow); ok {
					var rowData []string
					for _, cell := range rowNode.GetChildren() {
						if cellNode, ok := cell.(*ast.TableCell); ok {
							var cellText strings.Builder
							collectTextFromNode(cellNode, &cellText)
							rowData = append(rowData, strings.TrimSpace(cellText.String()))
						}
					}
					table.Rows = append(table.Rows, rowData)
				}
			}
		}
	}
	
	// Calculate column widths
	numCols := len(table.Headers)
	if numCols == 0 && len(table.Rows) > 0 {
		numCols = len(table.Rows[0])
	}
	
	table.ColWidths = make([]int, numCols)
	
	// Check header widths
	for i, header := range table.Headers {
		if i < len(table.ColWidths) && len(header) > table.ColWidths[i] {
			table.ColWidths[i] = len(header)
		}
	}
	
	// Check row widths
	for _, row := range table.Rows {
		for i, cell := range row {
			if i < len(table.ColWidths) && len(cell) > table.ColWidths[i] {
				table.ColWidths[i] = len(cell)
			}
		}
	}
	
	return table
}

// collectTextFromNode recursively collects text content from a node
func collectTextFromNode(node ast.Node, sb *strings.Builder) {
	switch n := node.(type) {
	case *ast.Text:
		sb.Write(n.Literal)
	default:
		for _, child := range n.GetChildren() {
			collectTextFromNode(child, sb)
		}
	}
}

// renderTable creates ASCII table representation
func renderTable(table *SimpleTable) string {
	if len(table.ColWidths) == 0 {
		return ""
	}
	
	var result strings.Builder
	
	// Render headers if present
	if len(table.Headers) > 0 {
		result.WriteString("| ")
		for i, header := range table.Headers {
			if i < len(table.ColWidths) {
				result.WriteString(header)
				result.WriteString(strings.Repeat(" ", table.ColWidths[i]-len(header)))
				result.WriteString(" | ")
			}
		}
		result.WriteString("\n")
		
		// Header separator
		result.WriteString("|")
		for _, width := range table.ColWidths {
			result.WriteString(strings.Repeat("-", width+2))
			result.WriteString("|")
		}
		result.WriteString("\n")
	}
	
	// Render data rows
	for _, row := range table.Rows {
		result.WriteString("| ")
		for i, cell := range row {
			if i < len(table.ColWidths) {
				result.WriteString(cell)
				result.WriteString(strings.Repeat(" ", table.ColWidths[i]-len(cell)))
				result.WriteString(" | ")
			}
		}
		result.WriteString("\n")
	}
	
	return result.String()
}

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
		sb.WriteString(`\sout{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Link:
		// For now, use footnote style: text\footnote{URL}
		processChildren(n)
		sb.WriteString(`\footnote{`)
		sb.WriteString(latex.EscapeLaTeX(string(n.Destination)))
		sb.WriteString(`}`)
	case *ast.Image:
		// For now, just render alt text (future: could include image)
		processChildren(n)
	case *ast.List:
		if n.ListFlags&ast.ListTypeOrdered != 0 {
			sb.WriteString(`\begin{enumerate}`)
			sb.WriteString("\n")
		} else {
			sb.WriteString(`\begin{itemize}`)
			sb.WriteString("\n")
		}
		
		for _, child := range n.GetChildren() {
			if listItem, ok := child.(*ast.ListItem); ok {
				sb.WriteString(`\item `)
				for _, itemChild := range listItem.GetChildren() {
					GenerateLatexFromAST(itemChild, sb)
				}
				sb.WriteString("\n")
			}
		}
		
		if n.ListFlags&ast.ListTypeOrdered != 0 {
			sb.WriteString(`\end{enumerate}`)
		} else {
			sb.WriteString(`\end{itemize}`)
		}
		sb.WriteString("\n")
	case *ast.ListItem:
		// This is handled by the List case above
		processChildren(n)
	case *ast.BlockQuote:
		sb.WriteString(`\begin{quote}`)
		sb.WriteString("\n")
		processChildren(n)
		sb.WriteString(`\end{quote}`)
		sb.WriteString("\n")
	case *ast.HorizontalRule:
		sb.WriteString(`\par\noindent\hrulefill\par`)
		sb.WriteString("\n")
	case *ast.Table:
		// Extract table data
		tableData := collectTableData(n)
		if len(tableData.ColWidths) == 0 {
			return
		}
		
		// Start tabular environment with left-aligned columns
		sb.WriteString(`\begin{tabular}{`)
		for range tableData.ColWidths {
			sb.WriteString(`l`)
		}
		sb.WriteString(`}`)
		sb.WriteString("\n")
		
		// Add headers if present
		if len(tableData.Headers) > 0 {
			for i, header := range tableData.Headers {
				if i > 0 {
					sb.WriteString(` & `)
				}
				sb.WriteString(latex.EscapeLaTeX(header))
			}
			sb.WriteString(` \\ \hline`)
			sb.WriteString("\n")
		}
		
		// Add data rows
		for _, row := range tableData.Rows {
			for i, cell := range row {
				if i > 0 {
					sb.WriteString(` & `)
				}
				if i < len(row) {
					sb.WriteString(latex.EscapeLaTeX(cell))
				}
			}
			sb.WriteString(` \\`)
			sb.WriteString("\n")
		}
		
		sb.WriteString(`\end{tabular}`)
		sb.WriteString("\n")
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
	case *ast.Del: // Handles ~~strikethrough~~
		sb.WriteString("\x1b[9m") // ANSI Strikethrough on
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[29m") // ANSI Strikethrough off
	case *ast.Link:
		sb.WriteString("\x1b[4m") // ANSI Underline on
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[24m") // ANSI Underline off
		// Append URL in dim color
		sb.WriteString("\x1b[2m (") // ANSI Dim on
		sb.WriteString(string(n.Destination))
		sb.WriteString(")\x1b[22m") // ANSI Dim off
	case *ast.Image:
		// Render alt text
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		// Append image reference in dim color
		sb.WriteString("\x1b[2m [img: ") // ANSI Dim on
		sb.WriteString(string(n.Destination))
		sb.WriteString("]\x1b[22m") // ANSI Dim off
	case *ast.List:
		// Process list items with proper numbering/bullets
		for i, child := range n.GetChildren() {
			if listItem, ok := child.(*ast.ListItem); ok {
				prefix := renderListPrefix(n, i)
				sb.WriteString(prefix)
				for _, itemChild := range listItem.GetChildren() {
					RenderMarkdownAST(itemChild, sb)
				}
				sb.WriteString("\n")
			}
		}
	case *ast.ListItem:
		// Process children directly (prefix is handled by List)
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
	case *ast.BlockQuote:
		// Process each line with blockquote prefix
		var content strings.Builder
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, &content)
		}
		// Split content by lines and add prefix
		lines := strings.Split(strings.TrimSuffix(content.String(), "\n"), "\n")
		prefix := renderBlockquotePrefix()
		for i, line := range lines {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(prefix)
			sb.WriteString(line)
		}
		sb.WriteString("\n")
	case *ast.HorizontalRule:
		sb.WriteString(renderHorizontalRule())
		sb.WriteString("\n")
	case *ast.Table:
		tableData := collectTableData(n)
		sb.WriteString(renderTable(tableData))
	case *ast.Code: // Inline code
		sb.WriteString("\x1b[7m") // ANSI Reverse video for inline code
		sb.WriteString(string(n.Literal))
		sb.WriteString("\x1b[27m") // ANSI Reverse video off
	case *ast.CodeBlock: // Code block
		sb.WriteString("\x1b[7m") // ANSI Reverse video for code block
		lines := strings.Split(string(n.Literal), "\n")
		for i, line := range lines {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(line)
		}
		sb.WriteString("\x1b[27m") // ANSI Reverse video off
		sb.WriteString("\n")
	case *ast.Heading:
		// Render heading with appropriate ANSI formatting
		level := n.Level
		switch level {
		case 1:
			sb.WriteString("\x1b[1;4m") // Bold + Underline
		case 2:
			sb.WriteString("\x1b[1m") // Bold
		case 3:
			sb.WriteString("\x1b[4m") // Underline
		default:
			sb.WriteString("\x1b[1m") // Bold for other levels
		}
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[0m") // Reset all formatting
		sb.WriteString("\n")
	// For common containers, just recurse on children
	case *ast.Document, *ast.Paragraph:
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
	case *ast.Math, *ast.MathBlock: // Include Math nodes in case parser produces them
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
	// Enhanced parser with common extensions for lists, tables, strikethrough, autolinks, etc.
	// Note: Exclude MathJax to avoid conflicts with the main math processing pipeline
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Strikethrough | parser.Tables | parser.Autolink
	// Remove MathJax from CommonExtensions if it's included
	extensions &^= parser.MathJax
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