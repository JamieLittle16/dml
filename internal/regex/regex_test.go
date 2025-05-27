package regex

import (
	"testing"
)

func TestInlineMath(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		content  string
	}{
		{"$E=mc^2$", true, "E=mc^2"},
		{"Some text $E=mc^2$ more text", true, "E=mc^2"},
		{"$$E=mc^2$$", false, ""},
		{"No math here", false, ""},
		{"$a$ $b$", true, "a"},
		{"$a\nb$", true, "a\nb"},
		{"$", false, ""},
		{"$$", false, ""},
	}

	for _, test := range tests {
		matches := InlineMath.FindStringSubmatch(test.input)
		found := len(matches) > 1
		if found != test.expected {
			t.Errorf("InlineMath.FindStringSubmatch(%q): got match = %v, want %v", test.input, found, test.expected)
			continue
		}

		if found && matches[1] != test.content {
			t.Errorf("InlineMath.FindStringSubmatch(%q): got content = %q, want %q", test.input, matches[1], test.content)
		}
	}
}

func TestDisplayMath(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		content  string
	}{
		{"$$E=mc^2$$", true, "E=mc^2"},
		{"Some text $$E=mc^2$$ more text", true, "E=mc^2"},
		{"$E=mc^2$", false, ""},
		{"No math here", false, ""},
		{"$$a$$ $$b$$", true, "a"},
		{"$$a\nb$$", true, "a\nb"},
		{"$$", false, ""},
	}

	for _, test := range tests {
		matches := DisplayMath.FindStringSubmatch(test.input)
		found := len(matches) > 1
		if found != test.expected {
			t.Errorf("DisplayMath.FindStringSubmatch(%q): got match = %v, want %v", test.input, found, test.expected)
			continue
		}

		if found && matches[1] != test.content {
			t.Errorf("DisplayMath.FindStringSubmatch(%q): got content = %q, want %q", test.input, matches[1], test.content)
		}
	}
}

func TestInlineMathParen(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		content  string
	}{
		{"\\(E=mc^2\\)", true, "E=mc^2"},
		{"Some text \\(E=mc^2\\) more text", true, "E=mc^2"},
		{"\\[E=mc^2\\]", false, ""},
		{"No math here", false, ""},
		{"\\(a\\) \\(b\\)", true, "a"},
		{"\\(a\nb\\)", true, "a\nb"},
		{"\\(", false, ""},
	}

	for _, test := range tests {
		matches := InlineMathParen.FindStringSubmatch(test.input)
		found := len(matches) > 1
		if found != test.expected {
			t.Errorf("InlineMathParen.FindStringSubmatch(%q): got match = %v, want %v", test.input, found, test.expected)
			continue
		}

		if found && matches[1] != test.content {
			t.Errorf("InlineMathParen.FindStringSubmatch(%q): got content = %q, want %q", test.input, matches[1], test.content)
		}
	}
}

func TestDisplayMathBracket(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		content  string
	}{
		{"\\[E=mc^2\\]", true, "E=mc^2"},
		{"Some text \\[E=mc^2\\] more text", true, "E=mc^2"},
		{"\\(E=mc^2\\)", false, ""},
		{"No math here", false, ""},
		{"\\[a\\] \\[b\\]", true, "a"},
		{"\\[a\nb\\]", true, "a\nb"},
		{"\\[", false, ""},
	}

	for _, test := range tests {
		matches := DisplayMathBracket.FindStringSubmatch(test.input)
		found := len(matches) > 1
		if found != test.expected {
			t.Errorf("DisplayMathBracket.FindStringSubmatch(%q): got match = %v, want %v", test.input, found, test.expected)
			continue
		}

		if found && matches[1] != test.content {
			t.Errorf("DisplayMathBracket.FindStringSubmatch(%q): got content = %q, want %q", test.input, matches[1], test.content)
		}
	}
}

func TestStreamProcessing(t *testing.T) {
	// Test start display math detection
	startTests := []struct {
		input    string
		expected bool
	}{
		{"$$formula", true},
		{" $$formula", true},
		{"Text $$formula", true},
		{"\\[formula", true},
		{"Text \\[formula", true},
		{"$formula$", false},
		{"\\(formula\\)", false},
		{"No math here", false},
	}

	for _, test := range startTests {
		found := StartDisplayMath.MatchString(test.input)
		if found != test.expected {
			t.Errorf("StartDisplayMath.MatchString(%q): got match = %v, want %v", test.input, found, test.expected)
		}
	}

	// Test end display math detection
	endTests := []struct {
		input    string
		expected bool
	}{
		{"formula$$", true},
		{"formula$$ ", true},
		{"formula$$ text", true},
		{"formula\\]", true},
		{"formula\\] text", true},
		{"$formula$", false},
		{"\\(formula\\)", false},
		{"No math here", false},
	}

	for _, test := range endTests {
		found := EndDisplayMath.MatchString(test.input)
		if found != test.expected {
			t.Errorf("EndDisplayMath.MatchString(%q): got match = %v, want %v", test.input, found, test.expected)
		}
	}
}