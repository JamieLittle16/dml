package terminal

import (
	"strings"
	"testing"
)

func TestKittyInline(t *testing.T) {
	// Create a simple PNG image for testing (1x1 pixel, black)
	samplePNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xDE, 0x00, 0x00, 0x00,
		0x0C, 0x49, 0x44, 0x41, 0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D, 0xB0, 0x00, 0x00, 0x00,
		0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	tests := []struct {
		name         string
		img          []byte
		isDisplayMath bool
		targetRows   int
		expectError  bool
		checkContent func(string) bool
	}{
		{
			name:         "Inline math with default rows",
			img:          samplePNG,
			isDisplayMath: false,
			targetRows:   0,
			expectError:  false,
			checkContent: func(s string) bool {
				return !strings.Contains(s, "\n") // Inline math should not have trailing newline
			},
		},
		{
			name:         "Display math with default rows",
			img:          samplePNG,
			isDisplayMath: true,
			targetRows:   0,
			expectError:  false,
			checkContent: func(s string) bool {
				return strings.HasSuffix(s, "\n") // Display math should have trailing newline
			},
		},
		{
			name:         "Custom height (2 rows)",
			img:          samplePNG,
			isDisplayMath: false,
			targetRows:   2,
			expectError:  false,
			checkContent: func(s string) bool {
				// Should contain height parameter
				return strings.Contains(s, "h=2")
			},
		},
		{
			name:         "Empty image",
			img:          []byte{},
			isDisplayMath: false,
			targetRows:   0,
			expectError:  true,
			checkContent: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := KittyInline(test.img, test.isDisplayMath, test.targetRows)
			
			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if test.checkContent != nil && !test.checkContent(result) {
				t.Errorf("Output didn't match expected format. Got: %q", result)
			}
			
			// Basic check that output starts with escape sequence for Kitty protocol
			if !strings.HasPrefix(result, "\x1b_G") {
				t.Errorf("Output doesn't start with Kitty protocol escape sequence. Got: %q", result)
			}
		})
	}
}

func TestDebugMode(t *testing.T) {
	// Test that setting debug mode works
	SetDebug(true)
	// Simply verify no panic occurs - actual debug output is to stderr
	SetDebug(false)
}

func TestNullCharacterHandling(t *testing.T) {
	// Create a simple PNG image for testing (1x1 pixel, black)
	samplePNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xDE, 0x00, 0x00, 0x00,
		0x0C, 0x49, 0x44, 0x41, 0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D, 0xB0, 0x00, 0x00, 0x00,
		0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	result, err := KittyInline(samplePNG, false, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that null characters are removed
	if strings.Contains(result, "\x00") {
		t.Errorf("Output contains null characters which should have been removed")
	}
}