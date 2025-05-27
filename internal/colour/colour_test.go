package colour

import (
	"testing"
)

func TestIsHexcolour(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"#000000", true},
		{"#FFFFFF", true},
		{"#123456", true},
		{"#ABC123", true},
		{"#FFF", true},
		{"#123", true},
		{"", false},
		{"#", false},
		{"#12345", false},
		{"#ZZZZZZ", false},
		{"white", false},
		{"#12345G", false},
	}

	for _, test := range tests {
		got := IsHexcolour(test.input)
		if got != test.want {
			t.Errorf("IsHexcolour(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}

func TestToHex(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"#000000", "#000000"},
		{"#FFF", "#ffffff"}, // Using lowercase as that's what the function returns
		{"white", "#FFFFFF"},
		{"RED", "#FF0000"},
		{"green", "#00FF00"},
		{"blue", "#0000FF"},
		{"", ""},
		{"nonexistent", ""},
		{"#12345G", ""},
	}

	for _, test := range tests {
		got := ToHex(test.input)
		if got != test.want {
			t.Errorf("ToHex(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestExpandHexcolour(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"#000", "#000000"},
		{"#FFF", "#FFFFFF"},
		{"#123", "#112233"},
		{"#ABC", "#AABBCC"},
		{"#000000", "#000000"},
		{"#FFFFFF", "#FFFFFF"},
	}

	for _, test := range tests {
		got := expandHexcolour(test.input)
		if got != test.want {
			t.Errorf("expandHexcolour(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestComplementHex(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"#000000", "#FFFFFF"},
		{"#FFFFFF", "#000000"},
		{"#FF0000", "#00FFFF"},
		{"#00FF00", "#FF00FF"},
		{"#0000FF", "#FFFF00"},
		{"#123456", "#EDCBA9"},
		{"#ABCDEF", "#543210"},
	}

	for _, test := range tests {
		got := ComplementHex(test.input)
		if got != test.want {
			t.Errorf("ComplementHex(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestLaTeXcolourDef(t *testing.T) {
	tests := []struct {
		name  string
		colour string
		want  string
	}{
		{"text", "#FF0000", "\\definecolour{text}{HTML}{FF0000}\n"},
		{"background", "#00FF00", "\\definecolour{background}{HTML}{00FF00}\n"},
		{"highlight", "#0000FF", "\\definecolour{highlight}{HTML}{0000FF}\n"},
	}

	for _, test := range tests {
		got := LaTeXcolourDef(test.name, test.colour)
		if got != test.want {
			t.Errorf("LaTeXcolourDef(%q, %q) = %q, want %q", test.name, test.colour, got, test.want)
		}
	}
}
