// Package terminal provides terminal-specific functionality for DML
package terminal

import (
	"bytes"
	"fmt"

	// "io"
	"os"
	"strings"

	"github.com/BourgeoisBear/rasterm"
)

var isDebug bool

// SetDebug enables or disables debug mode
func SetDebug(debug bool) {
	isDebug = debug
}

// KittyInline generates the Kitty graphics protocol string for the given image bytes
func KittyInline(img []byte, isDisplayMath bool, userTargetRows int) (string, error) {
	var sb strings.Builder
	opts := rasterm.KittyImgOpts{}

	// Set size based on math type and user preferences
	if userTargetRows > 0 {
		// User specified a size - honor it exactly
		opts.DstRows = uint32(userTargetRows)
	} else {
		// Use sensible defaults
		if isDisplayMath {
			// For display math, let Kitty decide the size
			// Auto-sizing (0) works well for display math
		} else {
			// For inline math, always use 1 row
			opts.DstRows = 1
		}
	}

	// Convert image to Kitty protocol
	err := rasterm.KittyCopyPNGInline(&sb, bytes.NewReader(img), opts)
	if err != nil {
		return "", fmt.Errorf("rasterm.KittyCopyPNGInline failed: %v", err)
	}
	kittyStr := sb.String()

	if isDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: Generated Kitty protocol with options: rows=%v, isDisplay=%v\n",
			opts.DstRows, isDisplayMath)
	}

	// Handle newlines for display math and inline math differently
	if isDisplayMath {
		// For display math, ensure there's exactly one trailing newline
		kittyStr = strings.TrimRight(kittyStr, "\n") + "\n"
	} else {
		// For inline math, ensure there are no trailing newlines
		kittyStr = strings.TrimRight(kittyStr, "\n")
	}

	// Remove null characters that might appear
	kittyStr = strings.ReplaceAll(kittyStr, "\x00", "")

	return kittyStr, nil
}
