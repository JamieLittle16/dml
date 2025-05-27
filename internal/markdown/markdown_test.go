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