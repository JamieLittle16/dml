// Package regex provides regex pattern functionality for DML
package regex

import (
	"regexp"
)

var (
	// Math expression patterns - use (?s) to make '.' match newlines
	InlineMath         = regexp.MustCompile(`(?s)\$(.+?)\$`)
	DisplayMath        = regexp.MustCompile(`(?s)\$\$(.+?)\$\$`)
	InlineMathParen    = regexp.MustCompile(`(?s)\\\((.+?)\\\)`)    // For \( ... \)
	DisplayMathBracket = regexp.MustCompile(`(?s)\\\[(.+?)\\\]`)    // For \[ ... \]

	// Streaming processing patterns - used to find the start and end of delimiters
	StartDisplayMath   = regexp.MustCompile(`(?:^|\s)\$\$|\\\[`)
	EndDisplayMath     = regexp.MustCompile(`\$\$(?:$|\s)|\\\]`)
)