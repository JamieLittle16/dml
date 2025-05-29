// Package colour provides colour management for DML
package colour

import (
	"fmt"
	"math"
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

// GetFuzzLevel returns an appropriate ImageMagick fuzz level percentage for a given hex color.
// It analyzes the color brightness and other characteristics to determine an optimal fuzz
// level for transparency detection.
func GetFuzzLevel(hexColor string) string {
	if hexColor == "" {
		return "50%" // Default for unrecognized colors
	}

	// Expand short hex if needed
	hexColor = expandHexcolour(hexColor)

	// Parse RGB components
	r, _ := strconv.ParseUint(hexColor[1:3], 16, 8)
	g, _ := strconv.ParseUint(hexColor[3:5], 16, 8)
	b, _ := strconv.ParseUint(hexColor[5:7], 16, 8)

	// Calculate perceived brightness
	// Using the formula: (0.299*R + 0.587*G + 0.114*B)
	brightness := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255.0

	// Calculate color saturation
	// Find min and max of RGB values
	max := math.Max(float64(r), math.Max(float64(g), float64(b)))
	min := math.Min(float64(r), math.Min(float64(g), float64(b)))
	saturation := 0.0
	if max > 0 {
		saturation = (max - min) / max
	}

	// Start with a more moderate base fuzz level
	var fuzzLevel float64 = 45.0
	
	// Special case for white and very bright colors with low saturation
	// White works well at 30%, so we'll keep it lower for very bright colors
	if brightness > 0.9 && saturation < 0.1 {
		fuzzLevel = 30.0
	}
	
	// Special handling for primary saturated colors (red and blue)
	// Check if this is primarily red (high R, low G and B)
	if float64(r) > 200 && float64(g) < 100 && float64(b) < 100 {
		fuzzLevel = 65.0 // Higher fuzz for reds (65%)
	}
	
	// Check if this is primarily blue (high B, low R and G)
	if float64(b) > 200 && float64(r) < 100 && float64(g) < 100 {
		fuzzLevel = 70.0 // Even higher fuzz for blues (70%)
	}
	
	// Check for other saturated primary/secondary colors and give them higher fuzz
	// But not as high as red/blue since they weren't specifically mentioned
	if saturation > 0.7 && brightness > 0.3 && !(
		// Skip the colors we've already handled (red and blue)
		(float64(r) > 200 && float64(g) < 100 && float64(b) < 100) || 
		(float64(b) > 200 && float64(r) < 100 && float64(g) < 100)) {
		fuzzLevel = 55.0
	}
	
	return fmt.Sprintf("%.1f%%", fuzzLevel)
}
