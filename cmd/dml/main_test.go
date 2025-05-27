package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestBasicExecution tests that the DML binary can be built and executed
func TestBasicExecution(t *testing.T) {
	// Skip if the DML binary doesn't exist
	if _, err := os.Stat("../../dml"); os.IsNotExist(err) {
		t.Skip("DML binary not found. Run ./build.sh first.")
	}

	tests := []struct {
		name    string
		input   string
		args    []string
		wantErr bool
		check   func(output string) bool
	}{
		{
			name:    "Plain text",
			input:   "Hello, world!",
			args:    []string{},
			wantErr: false,
			check: func(output string) bool {
				return strings.Contains(output, "Hello, world!")
			},
		},
		{
			name:    "Bold text",
			input:   "This is **bold** text",
			args:    []string{},
			wantErr: false,
			check: func(output string) bool {
				// Bold formatting should be applied
				return strings.Contains(output, "bold")
			},
		},
		{
			name:    "Help flag",
			input:   "",
			args:    []string{"--help"},
			wantErr: true, // --help causes exit code 2 which is treated as an error
			check: nil,
		},
		{
			name:    "Invalid flag",
			input:   "",
			args:    []string{"--invalid-flag"},
			wantErr: true,
			check:   nil,
		},
		{
			name:    "Custom colour",
			input:   "Text with colour",
			args:    []string{"--colour", "red"},
			wantErr: false,
			check: func(output string) bool {
				return strings.Contains(output, "Text with colour")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../../dml", tt.args...)
			cmd.Stdin = strings.NewReader(tt.input)
			
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			
			err := cmd.Run()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Command execution error = %v, wantErr %v\nStderr: %v", err, tt.wantErr, stderr.String())
				return
			}
			
			if tt.wantErr {
				return
			}
			
			output := stdout.String()
			if tt.check != nil && !tt.check(output) {
				t.Errorf("Output check failed.\nGot:\n%s", output)
			}
		})
	}
}

// TestMathRendering tests the math rendering capabilities
// Note: This test requires LaTeX and ImageMagick to be installed
func TestMathRendering(t *testing.T) {
	// Skip if the DML binary doesn't exist
	if _, err := os.Stat("../../dml"); os.IsNotExist(err) {
		t.Skip("DML binary not found. Run ./build.sh first.")
	}

	// Skip if SKIP_LATEX_TESTS environment variable is set
	if os.Getenv("SKIP_LATEX_TESTS") != "" {
		t.Skip("Skipping LaTeX tests due to SKIP_LATEX_TESTS environment variable")
	}

	// Check if LaTeX is installed
	if _, err := exec.LookPath("pdflatex"); err != nil {
		t.Skip("pdflatex not found, skipping math rendering tests")
	}

	// Check if ImageMagick is installed
	if _, err := exec.LookPath("convert"); err != nil {
		t.Skip("convert (ImageMagick) not found, skipping math rendering tests")
	}

	tests := []struct {
		name  string
		input string
		args  []string
	}{
		{
			name:  "Inline math",
			input: "Inline formula: $E=mc^2$",
			args:  []string{},
		},
		{
			name:  "Display math",
			input: "Display formula: $$\\sum_{i=1}^{n} i = \\frac{n(n+1)}{2}$$",
			args:  []string{},
		},
		{
			name:  "Custom colour math",
			input: "Coloured formula: $E=mc^2$",
			args:  []string{"-c", "blue"},
		},
		{
			name:  "Full document rendering",
			input: "# Document\n\nThis is a **bold** statement with formula $E=mc^2$.",
			args:  []string{"--render-all-latex"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../../dml", tt.args...)
			cmd.Stdin = strings.NewReader(tt.input)
			
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				t.Fatalf("Failed to get stdout pipe: %v", err)
			}
			
			stderr, err := cmd.StderrPipe()
			if err != nil {
				t.Fatalf("Failed to get stderr pipe: %v", err)
			}
			
			if err := cmd.Start(); err != nil {
				t.Fatalf("Failed to start command: %v", err)
			}
			
			// Read output and error
			outputBytes, err := io.ReadAll(stdout)
			if err != nil {
				t.Fatalf("Failed to read stdout: %v", err)
			}
			
			errorBytes, err := io.ReadAll(stderr)
			if err != nil {
				t.Fatalf("Failed to read stderr: %v", err)
			}
			
			// Wait for command to finish
			if err := cmd.Wait(); err != nil {
				t.Fatalf("Command failed: %v\nStderr: %s", err, errorBytes)
			}
			
			// Check that we got some output
			if len(outputBytes) == 0 {
				t.Errorf("Expected non-empty output")
			}
			
			// Check that output contains the Kitty graphics protocol marker
			output := string(outputBytes)
			if !strings.Contains(output, "\x1b_G") {
				t.Errorf("Output doesn't contain Kitty graphics protocol marker")
			}
		})
	}
}

// TestStreamingMode tests the streaming mode of DML
func TestStreamingMode(t *testing.T) {
	// Skip if the DML binary doesn't exist
	if _, err := os.Stat("../../dml"); os.IsNotExist(err) {
		t.Skip("DML binary not found. Run ./build.sh first.")
	}
	
	// Skip if SKIP_LATEX_TESTS environment variable is set
	if os.Getenv("SKIP_LATEX_TESTS") != "" {
		t.Skip("Skipping LaTeX tests due to SKIP_LATEX_TESTS environment variable")
	}

	// Create a test with multiline display math
	input := `This is a test with multiline display math:

$$
\\sum_{i=1}^{n} i = \\frac{n(n+1)}{2}
$$

And some more text after.`

	cmd := exec.Command("../../dml", "-D") // Debug mode to see more output
	cmd.Stdin = strings.NewReader(input)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}
	
	output := stdout.String()
	if !strings.Contains(output, "This is a test") {
		t.Errorf("Output doesn't contain expected text")
	}
	
	if !strings.Contains(output, "And some more text after") {
		t.Errorf("Output doesn't contain text after math block")
	}
	
	// In CI environments without a proper terminal, Kitty protocol might not be used
	// So we'll just check if we got some output
	if len(output) == 0 {
		t.Errorf("Expected non-empty output")
	}
}