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
	ListDepth        int
	OrderedIndex     int
	InBlockquote     bool
	BlockquotePrefix string
}

// getTerminalWidth tries to get terminal width, falls back to 80 chars
func getTerminalWidth() int {
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}
	return 80
}

// renderHorizontalRule creates a line of box-drawing characters for terminal display
func renderHorizontalRule() string {
	width := getTerminalWidth()
	if width > 40 {
		width = 40
	}
	return strings.Repeat("─", width)
}

// renderListPrefix returns the appropriate prefix for list items
func renderListPrefix(isList *ast.List, itemIndex int) string {
	if isList.ListFlags&ast.ListTypeOrdered != 0 {
		return strconv.Itoa(itemIndex+1) + ". "
	}
	return "• "
}

// renderBlockquotePrefix returns the prefix for blockquote lines
func renderBlockquotePrefix() string {
	return "│ "
}

// SimpleTable represents a basic table structure for ASCII rendering
type SimpleTable struct {
	Headers   []string
	Rows      [][]string
	ColWidths []int
}

// collectTableData extracts table data from AST nodes
func collectTableData(tableNode *ast.Table) *SimpleTable {
	table := &SimpleTable{}

	for _, child := range tableNode.GetChildren() {
		switch node := child.(type) {
		case *ast.TableHeader:
			for _, headerChild := range node.GetChildren() {
				if rowNode, ok := headerChild.(*ast.TableRow); ok {
					for _, cell := range rowNode.GetChildren() {
						if cellNode, ok := cell.(*ast.TableCell); ok {
							var cellText strings.Builder
							collectTextFromNode(cellNode, &cellText)
							table.Headers = append(table.Headers, strings.TrimSpace(cellText.String()))
						}
					}
					break
				}
			}
		case *ast.TableBody:
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

	numCols := len(table.Headers)
	if numCols == 0 && len(table.Rows) > 0 {
		numCols = len(table.Rows[0])
	}
	table.ColWidths = make([]int, numCols)

	for i, header := range table.Headers {
		if i < len(table.ColWidths) && len(header) > table.ColWidths[i] {
			table.ColWidths[i] = len(header)
		}
	}
	for _, row := range table.Rows {
		for i, cell := range row {
			if i < len(table.ColWidths) && len(cell) > table.ColWidths[i] {
				table.ColWidths[i] = len(cell)
			}
		}
	}

	return table
}

// collectTextFromNode recursively collects plain text content from a node
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

// renderTable creates a box-drawing table string using Unicode line characters.
//
// Example output:
//
//	┌───────┬───────┐
//	│ Name  │ Value │
//	├───────┼───────┤
//	│ alpha │ 1.0   │
//	│ beta  │ 2.0   │
//	└───────┴───────┘
func renderTable(table *SimpleTable) string {
	if len(table.ColWidths) == 0 {
		return ""
	}

	// Helpers to draw horizontal rule segments.
	hLine := func(left, mid, right, fill string) string {
		var b strings.Builder
		b.WriteString(left)
		for i, w := range table.ColWidths {
			b.WriteString(strings.Repeat(fill, w+2))
			if i < len(table.ColWidths)-1 {
				b.WriteString(mid)
			}
		}
		b.WriteString(right)
		b.WriteString("\n")
		return b.String()
	}

	cell := func(text string, width int) string {
		return " " + text + strings.Repeat(" ", width-len(text)) + " "
	}

	row := func(cells []string, sep string) string {
		var b strings.Builder
		b.WriteString(sep)
		for i, c := range cells {
			if i < len(table.ColWidths) {
				b.WriteString(cell(c, table.ColWidths[i]))
				b.WriteString(sep)
			}
		}
		b.WriteString("\n")
		return b.String()
	}

	var result strings.Builder

	// Top border
	result.WriteString(hLine("┌", "┬", "┐", "─"))

	if len(table.Headers) > 0 {
		// Header row
		result.WriteString(row(table.Headers, "│"))
		// Header/body separator
		result.WriteString(hLine("├", "┼", "┤", "─"))
	}

	for _, r := range table.Rows {
		result.WriteString(row(r, "│"))
	}

	// Bottom border
	result.WriteString(hLine("└", "┴", "┘", "─"))

	return result.String()
}

// GenerateLatexFromAST traverses a markdown AST and generates LaTeX output
func GenerateLatexFromAST(node ast.Node, sb *strings.Builder) {
	if node == nil {
		return
	}

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
	case *ast.Emph:
		sb.WriteString(`\textit{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Strong:
		sb.WriteString(`\textbf{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Del: // ~~strikethrough~~
		sb.WriteString(`\sout{`)
		processChildren(n)
		sb.WriteString(`}`)
	case *ast.Link:
		processChildren(n)
		sb.WriteString(`\footnote{`)
		sb.WriteString(latex.EscapeLaTeX(string(n.Destination)))
		sb.WriteString(`}`)
	case *ast.Image:
		// Render alt text only (no actual image embedding)
		processChildren(n)
	case *ast.List:
		if n.ListFlags&ast.ListTypeOrdered != 0 {
			sb.WriteString("\\begin{enumerate}\n")
		} else {
			sb.WriteString("\\begin{itemize}\n")
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
			sb.WriteString("\\end{enumerate}\n")
		} else {
			sb.WriteString("\\end{itemize}\n")
		}
	case *ast.ListItem:
		processChildren(n)
	case *ast.BlockQuote:
		sb.WriteString("\\begin{quote}\n")
		processChildren(n)
		sb.WriteString("\\end{quote}\n")
	case *ast.HorizontalRule:
		sb.WriteString("\\par\\noindent\\hrulefill\\par\n")
	case *ast.Table:
		tableData := collectTableData(n)
		if len(tableData.ColWidths) == 0 {
			return
		}
		sb.WriteString(`\begin{tabular}{`)
		for range tableData.ColWidths {
			sb.WriteString(`l`)
		}
		sb.WriteString("}\n")
		if len(tableData.Headers) > 0 {
			for i, header := range tableData.Headers {
				if i > 0 {
					sb.WriteString(` & `)
				}
				sb.WriteString(latex.EscapeLaTeX(header))
			}
			sb.WriteString(" \\\\ \\hline\n")
		}
		for _, row := range tableData.Rows {
			for i, cell := range row {
				if i > 0 {
					sb.WriteString(` & `)
				}
				if i < len(row) {
					sb.WriteString(latex.EscapeLaTeX(cell))
				}
			}
			sb.WriteString(" \\\\\n")
		}
		sb.WriteString("\\end{tabular}\n")
	case *ast.Math:
		sb.WriteString(`$`)
		sb.WriteString(string(n.Literal))
		sb.WriteString(`$`)
	case *ast.MathBlock:
		sb.WriteString(`$$`)
		sb.WriteString(string(n.Literal))
		sb.WriteString(`$$`)
	case *ast.Paragraph:
		processChildren(n)
		sb.WriteString("\n\\par\n\n")
	case *ast.Softbreak:
		sb.WriteString("\\\\\n")
	case *ast.Hardbreak:
		sb.WriteString("\\\\\n")
	case *ast.Code:
		sb.WriteString(`\texttt{`)
		sb.WriteString(latex.EscapeLaTeX(string(n.Literal)))
		sb.WriteString(`}`)
	case *ast.CodeBlock:
		sb.WriteString("\n\\begin{verbatim}\n")
		sb.WriteString(string(n.Literal))
		sb.WriteString("\\end{verbatim}\n")
	case *ast.Heading:
		switch n.Level {
		case 1:
			sb.WriteString(`\section*{`)
		case 2:
			sb.WriteString(`\subsection*{`)
		case 3:
			sb.WriteString(`\subsubsection*{`)
		default:
			sb.WriteString(`\paragraph*{`)
		}
		processChildren(n)
		sb.WriteString("}\n")
	default:
		for _, child := range n.GetChildren() {
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
	case *ast.Emph:
		sb.WriteString("\x1b[3m")
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[23m")
	case *ast.Strong:
		sb.WriteString("\x1b[1m")
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[22m")
	case *ast.Del: // ~~strikethrough~~
		sb.WriteString("\x1b[9m")
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[29m")
	case *ast.Link:
		sb.WriteString("\x1b[4m") // Underline
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[24m")
		sb.WriteString("\x1b[2m (")
		sb.WriteString(string(n.Destination))
		sb.WriteString(")\x1b[22m")
	case *ast.Image:
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[2m [img: ")
		sb.WriteString(string(n.Destination))
		sb.WriteString("]\x1b[22m")
	case *ast.List:
		for i, child := range n.GetChildren() {
			if listItem, ok := child.(*ast.ListItem); ok {
				sb.WriteString(renderListPrefix(n, i))
				for _, itemChild := range listItem.GetChildren() {
					RenderMarkdownAST(itemChild, sb)
				}
				sb.WriteString("\n")
			}
		}
	case *ast.ListItem:
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
	case *ast.BlockQuote:
		var content strings.Builder
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, &content)
		}
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
	case *ast.Code:
		sb.WriteString("\x1b[7m") // Reverse video
		sb.WriteString(string(n.Literal))
		sb.WriteString("\x1b[27m")
	case *ast.CodeBlock:
		sb.WriteString("\x1b[7m")
		lines := strings.Split(string(n.Literal), "\n")
		for i, line := range lines {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(line)
		}
		sb.WriteString("\x1b[27m\n")
	case *ast.Heading:
		switch n.Level {
		case 1:
			sb.WriteString("\x1b[1;4m") // Bold + underline
		case 2:
			sb.WriteString("\x1b[1m") // Bold
		case 3:
			sb.WriteString("\x1b[4m") // Underline
		default:
			sb.WriteString("\x1b[1m")
		}
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
		sb.WriteString("\x1b[0m\n")
	case *ast.Document, *ast.Paragraph:
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
	case *ast.Math, *ast.MathBlock:
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
	default:
		for _, child := range n.GetChildren() {
			RenderMarkdownAST(child, sb)
		}
	}
}

// ApplyFormatting parses a line (or block) of Markdown and converts it to
// ANSI-formatted text for terminal display.
func ApplyFormatting(line string) string {
	// Enable all common extensions except MathJax (handled separately by main.go)
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Strikethrough | parser.Tables | parser.Autolink
	extensions &^= parser.MathJax
	p := parser.NewWithExtensions(extensions)
	docNode := p.Parse([]byte(line))

	var sb strings.Builder
	RenderMarkdownAST(docNode, &sb)

	result := sb.String()
	result = strings.ReplaceAll(result, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\%", "%")
	result = strings.TrimSuffix(result, "%")
	result = strings.TrimSuffix(result, "\x00")

	return result
}
