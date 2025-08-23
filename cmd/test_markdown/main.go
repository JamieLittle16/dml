package main

import (
	"fmt"
	"strings"
	"dml/internal/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func testTable() {
	input := `| Col A | Col B |
|-------|-------|
| 1     | 2     |
| 3     | 4     |`
	
	fmt.Printf("=== Table Test ===\n")
	fmt.Printf("Input:\n%s\n\n", input)
	
	// Test ANSI rendering
	result := markdown.ApplyFormatting(input)
	fmt.Printf("ANSI output:\n%s\n", result)
	
	// Test LaTeX generation
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Strikethrough | parser.Tables | parser.Autolink
	extensions &^= parser.MathJax
	p := parser.NewWithExtensions(extensions)
	docNode := p.Parse([]byte(input))
	
	var sb strings.Builder
	markdown.GenerateLatexFromAST(docNode, &sb)
	latexResult := sb.String()
	
	fmt.Printf("LaTeX output:\n%s\n", latexResult)
}

func main() {
	testTable()
}