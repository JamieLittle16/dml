// Package latex provides LaTeX rendering functionality for DML
package latex

import (
	"strings"
)

// Special characters that need escaping in LaTeX
// Order matters for some replacements (e.g., `\` before other chars that might use it)
var latexEscaper = strings.NewReplacer(
	`\`, `\textbackslash{}`, // Must be first
	`&`, `\&`,
	`%`, `\%`,
	`$`, `\$`,
	`#`, `\#`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
	`~`, `\textasciitilde{}`,
	`^`, `\textasciicircum{}`,
	`[`, `{[}`,
	`]`, `{]}`,
	`|`, `{\vert}`,
	`/`, `{/}`,
)

// EscapeLaTeX escapes special characters in a string for LaTeX
func EscapeLaTeX(s string) string {
	return latexEscaper.Replace(s)
}