// Package colour provides colour management for DML
package colour

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// IsHexcolour returns true if the string is a valid hex colour (#RRGGBB or #RGB)
func IsHexcolour(s string) bool {
	matched, _ := regexp.MatchString(`^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$`, s)
	return matched
}

// Map of named colours to hex codes (CSS/X11 names)
var namedcolourHex = map[string]string{
	"black":   "#000000",
	"white":   "#FFFFFF",
	"red":     "#FF0000",
	"green":   "#00FF00",
	"blue":    "#0000FF",
	"yellow":  "#FFFF00",
	"cyan":    "#00FFFF",
	"magenta": "#FF00FF",
	"gray":    "#808080",
	"grey":    "#808080",
	"orange":  "#FFA500",
	"purple":  "#800080",
	"brown":   "#A52A2A",
	"pink":    "#FFC0CB",
	"lime":    "#00FF00",
	"navy":    "#000080",
	"teal":    "#008080",
	"maroon":  "#800000",
	"olive":   "#808000",
	"silver":  "#C0C0C0",
}

// ToHex returns a hex code for a colour string (hex or named), or empty string if unknown
func ToHex(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if IsHexcolour(s) {
		return expandHexcolour(s)
	}
	if hex, ok := namedcolourHex[s]; ok {
		return hex
	}
	return ""
}

// expandHexcolour expands #RGB to #RRGGBB
func expandHexcolour(s string) string {
	if len(s) == 4 {
		return "#" + strings.Repeat(string(s[1]), 2) +
			strings.Repeat(string(s[2]), 2) +
			strings.Repeat(string(s[3]), 2)
	}
	return s
}

// ComplementHex returns the complement of a #RRGGBB hex colour as #RRGGBB
func ComplementHex(s string) string {
	s = expandHexcolour(s)
	r, _ := strconv.ParseUint(s[1:3], 16, 8)
	g, _ := strconv.ParseUint(s[3:5], 16, 8)
	b, _ := strconv.ParseUint(s[5:7], 16, 8)
	return fmt.Sprintf("#%02X%02X%02X", 0xFF^r, 0xFF^g, 0xFF^b)
}

// LaTeXcolourDef returns the colour name and LaTeX colour definition for a hex colour
func LaTeXcolourDef(name, hex string) string {
	return fmt.Sprintf("\\definecolor{%s}{HTML}{%s}\n", name, hex[1:])
}
