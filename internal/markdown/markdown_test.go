package markdown

import (
	"strings"
	"testing"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func TestRenderMarkdownAST(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "Plain text",
			markdown: "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "Bold text",
			markdown: "This is **bold** text",
			expected: "This is \x1b[1mbold\x1b[22m text",
		},
		{
			name:     "Italic text",
			markdown: "This is *italic* text",
			expected: "This is \x1b[3mitalic\x1b[23m text",
		},
		{
			name:     "Mixed formatting",
			markdown: "This is **bold** and *italic* text",
			expected: "This is \x1b[1mbold\x1b[22m and \x1b[3mitalic\x1b[23m text",
		},
		{
			name:     "Nested formatting",
			markdown: "This is ***bold and italic*** text",
			expected: "This is \x1b[1m\x1b[3mbold and italic\x1b[23m\x1b[22m text",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := parser.New()
			doc := p.Parse([]byte(test.markdown))

			var sb strings.Builder
			RenderMarkdownAST(doc, &sb)

			got := sb.String()
			if got != test.expected {
				t.Errorf("Expected: %q, got: %q", test.expected, got)
			}
		})
	}
}

func TestApplyFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Plain text",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "Bold text",
			input:    "This is **bold** text",
			expected: "This is \x1b[1mbold\x1b[22m text",
		},
		{
			name:     "Italic text",
			input:    "This is *italic* text",
			expected: "This is \x1b[3mitalic\x1b[23m text",
		},
		{
			name:     "Text with newlines",
			input:    "Line 1\nLine 2",
			expected: "Line 1\nLine 2",
		},
		{
			name:     "Text with escaped characters",
			input:    "Escaped \\n newline and \\% percent",
			expected: "Escaped \n newline and % percent",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ApplyFormatting(test.input)
			if got != test.expected {
				t.Errorf("Expected: %q, got: %q", test.expected, got)
			}
		})
	}
}

func TestGenerateLatexFromAST(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "Plain text",
			markdown: "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "Bold text",
			markdown: "This is **bold** text",
			expected: "This is \\textbf{bold} text",
		},
		{
			name:     "Italic text",
			markdown: "This is *italic* text",
			expected: "This is \\textit{italic} text",
		},
		{
			name:     "Special characters",
			markdown: "Symbols: & % $ # _ { } ~ ^ \\ [ ]",
			expected: "Symbols: \\& \\% \\$ \\# \\_ \\{ \\} \\textasciitilde{} \\textasciicircum{} \\textbackslash{} {[} {]}",
		},
		{
			name:     "Paragraph",
			markdown: "Paragraph 1\n\nParagraph 2",
			expected: "Paragraph 1\n\\par\n\nParagraph 2\n\\par\n\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := parser.New()
			doc := p.Parse([]byte(test.markdown))

			var sb strings.Builder
			GenerateLatexFromAST(doc, &sb)

			got := sb.String()
			if !strings.Contains(got, test.expected) {
				t.Errorf("Expected result to contain: %q, got: %q", test.expected, got)
			}
		})
	}
}

// Mock AST node for testing
type mockNode struct {
	ast.Node
	children []ast.Node
	literal  []byte
}

func (m *mockNode) GetChildren() []ast.Node {
	return m.children
}

func TestEmptyNode(t *testing.T) {
	var sb strings.Builder
	
	// Test with nil node
	RenderMarkdownAST(nil, &sb)
	if sb.Len() != 0 {
		t.Errorf("Expected empty string for nil node, got: %q", sb.String())
	}
	
	sb.Reset()
	GenerateLatexFromAST(nil, &sb)
	if sb.Len() != 0 {
		t.Errorf("Expected empty string for nil node, got: %q", sb.String())
	}
}

// Test enhanced features with parser extensions
func TestEnhancedMarkdownFeatures(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expectedANSI string
		expectedLaTeX string
	}{
		{
			name:     "Strikethrough",
			markdown: "This is ~~strikethrough~~ text",
			expectedANSI: "This is \x1b[9mstrikethrough\x1b[29m text",
			expectedLaTeX: "This is \\sout{strikethrough} text",
		},
		{
			name:     "Link",
			markdown: "A [link](https://example.com) here",
			expectedANSI: "A \x1b[4mlink\x1b[24m\x1b[2m (https://example.com)\x1b[22m here",
			expectedLaTeX: "A link\\footnote{https:{/}{/}example.com} here",
		},
		{
			name:     "Image",
			markdown: "An ![image](test.jpg) here",
			expectedANSI: "An image\x1b[2m [img: test.jpg]\x1b[22m here",
			expectedLaTeX: "An image here",
		},
		{
			name:     "Ordered list",
			markdown: "1. First item\n2. Second item",
			expectedANSI: "1. First item\n2. Second item\n",
			expectedLaTeX: "\\begin{enumerate}\n\\item First item\n\\par\n\n\n\\item Second item\n\\par\n\n\n\\end{enumerate}\n",
		},
		{
			name:     "Unordered list",
			markdown: "- First item\n- Second item",
			expectedANSI: "• First item\n• Second item\n",
			expectedLaTeX: "\\begin{itemize}\n\\item First item\n\\par\n\n\n\\item Second item\n\\par\n\n\n\\end{itemize}\n",
		},
		{
			name:     "Blockquote",
			markdown: "> This is a quote",
			expectedANSI: "│ This is a quote\n",
			expectedLaTeX: "\\begin{quote}\nThis is a quote\n\\par\n\n\\end{quote}\n",
		},
		{
			name:     "Horizontal rule",
			markdown: "---",
			expectedANSI: "────────────────────────────────────────\n",
			expectedLaTeX: "\\par\\noindent\\hrulefill\\par\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name+" ANSI", func(t *testing.T) {
			result := ApplyFormatting(test.markdown)
			if result != test.expectedANSI {
				t.Errorf("ANSI: Expected %q, got %q", test.expectedANSI, result)
			}
		})

		t.Run(test.name+" LaTeX", func(t *testing.T) {
			// Parse with enhanced extensions
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Strikethrough | parser.Tables | parser.Autolink
			extensions &^= parser.MathJax
			p := parser.NewWithExtensions(extensions)
			doc := p.Parse([]byte(test.markdown))

			var sb strings.Builder
			GenerateLatexFromAST(doc, &sb)

			result := sb.String()
			if !strings.Contains(result, test.expectedLaTeX) {
				t.Errorf("LaTeX: Expected result to contain %q, got %q", test.expectedLaTeX, result)
			}
		})
	}
}

// Test table functionality specifically
func TestTableRendering(t *testing.T) {
	tableMarkdown := `| Col A | Col B |
|-------|-------|
| 1     | 2     |
| 3     | 4     |`

	// Test ANSI rendering
	ansiResult := ApplyFormatting(tableMarkdown)
	
	// Check that we get a table structure (exact format may vary)
	if !strings.Contains(ansiResult, "Col A") || !strings.Contains(ansiResult, "Col B") {
		t.Errorf("ANSI table should contain headers, got: %q", ansiResult)
	}
	if !strings.Contains(ansiResult, "| 1") || !strings.Contains(ansiResult, "| 3") {
		t.Errorf("ANSI table should contain data, got: %q", ansiResult)
	}

	// Test LaTeX rendering
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Strikethrough | parser.Tables | parser.Autolink
	extensions &^= parser.MathJax
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(tableMarkdown))

	var sb strings.Builder
	GenerateLatexFromAST(doc, &sb)
	latexResult := sb.String()

	expectedLaTeX := []string{
		"\\begin{tabular}{ll}",
		"Col A & Col B",
		"\\hline",
		"1 & 2",
		"3 & 4",
		"\\end{tabular}",
	}

	for _, expected := range expectedLaTeX {
		if !strings.Contains(latexResult, expected) {
			t.Errorf("LaTeX table should contain %q, got: %q", expected, latexResult)
		}
	}
}

// Test mixed formatting
func TestMixedFormatting(t *testing.T) {
	markdown := "This is **bold** and ~~strikethrough~~ with a [link](https://example.com)."
	result := ApplyFormatting(markdown)
	
	// Should contain ANSI codes for bold, strikethrough, and link
	if !strings.Contains(result, "\x1b[1m") {
		t.Error("Should contain bold ANSI code")
	}
	if !strings.Contains(result, "\x1b[9m") {
		t.Error("Should contain strikethrough ANSI code")
	}
	if !strings.Contains(result, "\x1b[4m") {
		t.Error("Should contain underline ANSI code for link")
	}
	if !strings.Contains(result, "https://example.com") {
		t.Error("Should contain link URL")
	}
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	// Test getTerminalWidth
	width := getTerminalWidth()
	if width <= 0 {
		t.Error("Terminal width should be positive")
	}

	// Test renderHorizontalRule
	rule := renderHorizontalRule()
	if len(rule) == 0 {
		t.Error("Horizontal rule should not be empty")
	}
	if !strings.Contains(rule, "─") && !strings.Contains(rule, "-") {
		t.Error("Horizontal rule should contain line characters")
	}

	// Test renderBlockquotePrefix
	prefix := renderBlockquotePrefix()
	if prefix != "│ " {
		t.Errorf("Expected blockquote prefix '│ ', got %q", prefix)
	}
}