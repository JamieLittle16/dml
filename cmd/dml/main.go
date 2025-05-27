// Package main is the entry point for the DML tool
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"dml/internal/latex"
	"dml/internal/markdown"
	"dml/internal/regex"
	"dml/internal/terminal"

	"github.com/gomarkdown/markdown/parser"
)

func main() {
	// Command-line flags
	colourFlag := flag.String("colour", "white", "Set LaTeX text colour (e.g., red, #00FF00).")
	cFlag := flag.String("c", "", "Short alias for --colour. Overrides --colour if set.")
	sizeFlag := flag.Int("size", 0, "Target terminal rows for LaTeX images (0 for default: 1 for inline, auto for display).")
	sFlag := flag.Int("s", 0, "Short alias for --size.")
	dpiFlag := flag.Int("dpi", 300, "Set DPI for rendering LaTeX images.")
	dFlag := flag.Int("d", 0, "Short alias for --dpi. Overrides --dpi if set (and not 0).")
	renderAllLatexFlag := flag.Bool("render-all-latex", false, "Render entire input as a single LaTeX document/image.")
	lFlag := flag.Bool("l", false, "Short alias for --render-all-latex.")
	debugFlag := flag.Bool("debug", false, "Enable verbose debug output.")
	dDebugFlag := flag.Bool("D", false, "Short alias for --debug.")

	flag.Parse() // Parse all flags first

	// Set default colour to white and apply overrides if specified
	effectivecolour := "white"
	if *cFlag != "" {
		effectivecolour = *cFlag
	} else if *colourFlag != "" && *colourFlag != "white" {
		effectivecolour = *colourFlag
	}

	// Determine if short flags -s and -d were explicitly set
	var sFlagSet, dFlagSet bool
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "s":
			sFlagSet = true
		case "d":
			dFlagSet = true
		}
	})

	// Determine the effective size to use
	effectiveSize := *sizeFlag
	if sFlagSet { // If -s was explicitly provided on the command line, it takes precedence
		effectiveSize = *sFlag
	}
	if effectiveSize < 0 { // Treat negative size as default (0)
		effectiveSize = 0
	}

	// Determine the effective DPI to use
	effectiveDPI := *dpiFlag
	if dFlagSet { // If -d was explicitly provided on the command line, it takes precedence
		effectiveDPI = *dFlag
	}
	if effectiveDPI <= 0 { // Ensure DPI is positive, default to 300 if invalid
		effectiveDPI = 300
	}

	isRenderAllLatexMode := *renderAllLatexFlag || *lFlag
	isDebugMode := *debugFlag || *dDebugFlag

	// Set debug mode in packages
	if isDebugMode {
		latex.SetDebug(true)
		terminal.SetDebug(true)
		os.Setenv("DML_DEBUG", "1")
		fmt.Fprintln(os.Stderr, "DEBUG: dml starting")
		fmt.Fprintln(os.Stderr, "DEBUG: Flags parsed.")
		fmt.Fprintf(os.Stderr, "DEBUG: isRenderAllLatexMode: %v\n", isRenderAllLatexMode)
	}

	if isRenderAllLatexMode {
		processFullDocument(effectivecolour, effectiveSize, effectiveDPI, isDebugMode)
	} else {
		processStreamingDocument(effectivecolour, effectiveSize, effectiveDPI, isDebugMode)
	}

	// Final debug messages if debug mode is enabled
	if isDebugMode {
		fmt.Fprintf(os.Stderr, "DEBUG: dml execution completed. If math rendering issues occurred, check for LaTeX or convert errors.")
		fmt.Fprintln(os.Stderr, "DEBUG: dml exiting.")
	}
}

// processFullDocument handles the full document rendering mode
func processFullDocument(effectivecolour string, effectiveSize, effectiveDPI int, isDebugMode bool) {
	if isDebugMode {
		fmt.Fprintln(os.Stderr, "DEBUG: Reading standard input (full document) for render-all-latex mode...")
	}

	// Read all of stdin into a single string
	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
		os.Exit(1)
	}
	inputString := string(inputBytes)

	if isDebugMode {
		fmt.Fprintf(os.Stderr, "DEBUG: Finished reading input (%d bytes).\n", len(inputBytes))
	}

	// Preprocess \[...\] and \(...\) to $$...$$ and $...$ for correct math parsing
	preprocessed := regex.DisplayMathBracket.ReplaceAllStringFunc(inputString, func(match string) string {
		content := strings.TrimSpace(match[2 : len(match)-2])
		return "$" + content + "$"
	})
	preprocessed = regex.InlineMathParen.ReplaceAllStringFunc(preprocessed, func(match string) string {
		content := strings.TrimSpace(match[2 : len(match)-2])
		return "$" + content + "$"
	})

	// Enable MathJax and other common extensions for parsing
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.MathJax)
	docNode := p.Parse([]byte(preprocessed))

	var latexBodyBuilder strings.Builder
	markdown.GenerateLatexFromAST(docNode, &latexBodyBuilder)
	latexBody := latexBodyBuilder.String()

	img, renderErr := latex.RenderFullDocument(latexBody, effectivecolour, effectiveDPI)
	if renderErr != nil {
		fmt.Fprintf(os.Stderr, "Error in full LaTeX rendering mode: %v\n", renderErr)
		// In full render mode, if LaTeX fails, print the original input so user can debug
		fmt.Print(inputString)
		os.Exit(1)
	}

	kittyStr, kittyErr := terminal.KittyInline(img, true, effectiveSize)
	if kittyErr != nil {
		fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for full document: %v\n", kittyErr)
		// If Kitty protocol generation fails, print original input
		fmt.Print(inputString)
		os.Exit(1)
	}
	fmt.Print(kittyStr)
}

// processStreamingDocument handles the streaming mode with line-by-line processing
func processStreamingDocument(effectivecolour string, effectiveSize, effectiveDPI int, isDebugMode bool) {
	if isDebugMode {
		fmt.Fprintln(os.Stderr, "DEBUG: Entering standard processing mode (line-by-line streaming with state).")
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout) // Use a buffered writer for output flushing

	var mathBuffer strings.Builder // Buffer for collecting multi-line math content
	inDisplayMath := false         // State flag

	if isDebugMode {
		fmt.Fprintln(os.Stderr, "DEBUG: Starting line-by-line input reading and processing loop...")
	}

	for {
		inputLine, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading input line: %v\n", err)
			writer.Flush() // Flush any pending output
			os.Exit(1)
		}

		// Determine if this is the last line
		isLastLine := (err == io.EOF)

		if inDisplayMath {
			// We are inside a display math block
			if isDebugMode {
				fmt.Fprintf(os.Stderr, "DEBUG: In display math mode, processing line: %s\n", strings.TrimSpace(inputLine))
			}

			endMatchIdx := regex.EndDisplayMath.FindStringIndex(inputLine)

			if endMatchIdx != nil {
				// Found the closing delimiter on this line
				if isDebugMode {
					fmt.Fprintln(os.Stderr, "DEBUG: Found display math end delimiter.")
				}
				mathBuffer.WriteString(inputLine[:endMatchIdx[0]]) // Add content before the delimiter

				mathContent := mathBuffer.String()
				mathBuffer.Reset() // Clear the buffer
				inDisplayMath = false // Exit display math state

				// Render the collected math content
				if isDebugMode {
					fmt.Fprintf(os.Stderr, "DEBUG: Attempting to render display math (length: %d chars)\n", len(mathContent))
				}
				img, renderErr := latex.RenderMath(mathContent, effectivecolour, true, effectiveDPI)
				if renderErr != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Rendering display math failed: %v\n", renderErr)
					// On error, print the un-rendered content as text
					writer.WriteString("$$")
					writer.WriteString(mathContent)
					writer.WriteString("$$\n") // Add newline if it's missing from content
				} else {
					if isDebugMode {
						fmt.Fprintln(os.Stderr, "DEBUG: Math rendering successful, generating Kitty protocol")
					}
					kittyStr, kittyErr := terminal.KittyInline(img, true, effectiveSize)
					if kittyErr != nil {
						fmt.Fprintf(os.Stderr, "ERROR: Generating Kitty protocol failed: %v\n", kittyErr)
						// On error, print the un-rendered content as text
						writer.WriteString("$$")
						writer.WriteString(mathContent)
						writer.WriteString("$$\n") // Add newline if it's missing from content
					} else {
						if isDebugMode {
							fmt.Fprintln(os.Stderr, "DEBUG: Successfully generated Kitty protocol for display math")
						}
						writer.WriteString(kittyStr) // Write the rendered image protocol
					}
				}

				// Process the rest of the line after the closing delimiter
				remainingLine := inputLine[endMatchIdx[1]:]
				if len(remainingLine) > 0 {
					// Process remaining part of the line as normal text
					if isDebugMode {
						fmt.Fprintf(os.Stderr, "DEBUG: Processing remaining line after display math: %s\n", strings.TrimSpace(remainingLine))
					}
					processedRemaining := processInlineMath(remainingLine, effectivecolour, effectiveSize, effectiveDPI, isDebugMode)
					finalRemainingOutput := markdown.ApplyFormatting(processedRemaining)
					writer.WriteString(finalRemainingOutput)
				}

			} else {
				// No closing delimiter yet, just buffer the line
				mathBuffer.WriteString(inputLine)
				if isDebugMode {
					fmt.Fprintln(os.Stderr, "DEBUG: Appended line to math buffer.")
				}
			}

		} else {
			// We are in normal text mode
			if isDebugMode {
				fmt.Fprintf(os.Stderr, "DEBUG: In normal mode, processing line: %s\n", strings.TrimSpace(inputLine))
			}

			startMatchIdx := regex.StartDisplayMath.FindStringIndex(inputLine)
			endMatchIdx := regex.EndDisplayMath.FindStringIndex(inputLine) // Check for same-line closing

			if startMatchIdx != nil && (endMatchIdx == nil || endMatchIdx[0] < startMatchIdx[0]) {
				// Found starting delimiter for a multi-line block (and no closing before it)
				if isDebugMode {
					fmt.Fprintln(os.Stderr, "DEBUG: Found display math start delimiter. Switching to math state.")
				}
				// Process content *before* the delimiter as normal text
				beforeDelimiter := inputLine[:startMatchIdx[0]]
				if len(beforeDelimiter) > 0 {
					if isDebugMode {
						fmt.Fprintf(os.Stderr, "DEBUG: Processing text before delimiter: %s\n", strings.TrimSpace(beforeDelimiter))
					}
					processedBefore := processInlineMath(beforeDelimiter, effectivecolour, effectiveSize, effectiveDPI, isDebugMode)
					finalBeforeOutput := markdown.ApplyFormatting(processedBefore)
					writer.WriteString(finalBeforeOutput)
				}

				// Start buffering from the content *after* the delimiter on this line
				mathBuffer.WriteString(inputLine[startMatchIdx[1]:])
				inDisplayMath = true // Enter display math state
				if isDebugMode {
					fmt.Fprintln(os.Stderr, "DEBUG: Started buffering math content.")
				}

			} else if startMatchIdx != nil && endMatchIdx != nil && startMatchIdx[0] < endMatchIdx[0] {
				// Found both start and end delimiters on the same line (single-line display math)
				if isDebugMode {
					fmt.Fprintln(os.Stderr, "DEBUG: Found single-line display math.")
				}
				// Process content *before* the start delimiter
				beforeDelimiter := inputLine[:startMatchIdx[0]]
				if len(beforeDelimiter) > 0 {
					if isDebugMode {
						fmt.Fprintf(os.Stderr, "DEBUG: Processing text before single-line math: %s\n", strings.TrimSpace(beforeDelimiter))
					}
					processedBefore := processInlineMath(beforeDelimiter, effectivecolour, effectiveSize, effectiveDPI, isDebugMode)
					finalBeforeOutput := markdown.ApplyFormatting(processedBefore)
					writer.WriteString(finalBeforeOutput)
				}

				// Extract and process the math content
				mathContent := inputLine[startMatchIdx[1]:endMatchIdx[0]]
				if isDebugMode {
					fmt.Fprintf(os.Stderr, "DEBUG: Rendering single-line display math: %s\n", strings.TrimSpace(mathContent))
				}
				img, renderErr := latex.RenderMath(mathContent, effectivecolour, true, effectiveDPI)
				if renderErr != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Rendering display math failed: %v\n", renderErr)
					// On error, print the un-rendered content as text
					writer.WriteString("$$")
					writer.WriteString(mathContent)
					writer.WriteString("$$\n") // Add newline if it's missing from content
				} else {
					kittyStr, kittyErr := terminal.KittyInline(img, true, effectiveSize)
					if kittyErr != nil {
						fmt.Fprintf(os.Stderr, "ERROR: Generating Kitty protocol failed: %v\n", kittyErr)
						// On error, print the un-rendered content as text
						writer.WriteString("$$")
						writer.WriteString(mathContent)
						writer.WriteString("$$\n") // Add newline if it's missing from content
					} else {
						writer.WriteString(kittyStr) // Write the rendered image protocol
					}
				}

				// Process content *after* the end delimiter
				afterDelimiter := inputLine[endMatchIdx[1]:]
				if len(afterDelimiter) > 0 {
					if isDebugMode {
						fmt.Fprintf(os.Stderr, "DEBUG: Processing text after single-line math: %s\n", strings.TrimSpace(afterDelimiter))
					}
					processedAfter := processInlineMath(afterDelimiter, effectivecolour, effectiveSize, effectiveDPI, isDebugMode)
					finalAfterOutput := markdown.ApplyFormatting(processedAfter)
					writer.WriteString(finalAfterOutput)
				}

			} else {
				// No display math delimiters found on this line.
				// Process for inline math and markdown as before.
				processedLine := processInlineMath(inputLine, effectivecolour, effectiveSize, effectiveDPI, isDebugMode)

				// Apply Markdown formatting to the processed line.
				finalLineOutput := markdown.ApplyFormatting(processedLine)

				// Remove any trailing special characters that might appear
				finalLineOutput = strings.TrimSuffix(finalLineOutput, "%")
				finalLineOutput = strings.TrimSuffix(finalLineOutput, "\x00")

				// Make sure we keep newlines as is
				if strings.HasSuffix(processedLine, "\n") && !strings.HasSuffix(finalLineOutput, "\n") {
					finalLineOutput += "\n"
				}

				// Write the processed line to the buffered writer.
				writer.WriteString(finalLineOutput)
			}
		}

		// Flush the writer immediately after processing a line (or block end).
		writer.Flush()

		// If the error was EOF, it means we just processed the last line.
		// The loop condition should break after this iteration.
		if isLastLine {
			if isDebugMode {
				fmt.Fprintln(os.Stderr, "DEBUG: Finished processing line before EOF. Exiting loop after flush.")
			}
			break
		}
	}

	// Handle case where input ended while still inside a display math block
	if inDisplayMath {
		if isDebugMode {
			fmt.Fprintln(os.Stderr, "DEBUG: Warning: Reached EOF while still inside a display math block. Outputting buffered content as plain text.")
		}
		writer.WriteString("$$") // Output the start delimiter that wasn't closed
		writer.WriteString(mathBuffer.String()) // Output the buffered content
		// No closing delimiter to output
		writer.Flush()
	}

	if isDebugMode {
		fmt.Fprintln(os.Stderr, "DEBUG: Output streaming finished.")
	}
}

// processInlineMath handles inline math expressions in a text line
func processInlineMath(line, effectivecolour string, effectiveSize, effectiveDPI int, isDebugMode bool) string {
	// Process $...$ inline math
	processedLine := regex.InlineMath.ReplaceAllStringFunc(line, func(match string) string {
		content := strings.TrimSpace(match[1 : len(match)-1])
		if content == "" { return match }

		if isDebugMode {
			fmt.Fprintf(os.Stderr, "DEBUG: Processing inline math with colour: '%s'\n", effectivecolour)
		}

		img, rErr := latex.RenderMath(content, effectivecolour, false, effectiveDPI)
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
			return match
		}
		kStr, kErr := terminal.KittyInline(img, false, effectiveSize)
		if kErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
			return match
		}
		return kStr
	})

	// Process \(...\) inline math
	processedLine = regex.InlineMathParen.ReplaceAllStringFunc(processedLine, func(match string) string {
		content := strings.TrimSpace(match[2 : len(match)-2])
		if content == "" { return match }

		if isDebugMode {
			fmt.Fprintf(os.Stderr, "DEBUG: Processing parenthesis-style inline math with colour: '%s'\n", effectivecolour)
		}

		img, rErr := latex.RenderMath(content, effectivecolour, false, effectiveDPI)
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "Error rendering inline math ('%s'): %v\n", content, rErr)
			return match
		}
		kStr, kErr := terminal.KittyInline(img, false, effectiveSize)
		if kErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating Kitty protocol for inline math ('%s'): %v\n", content, kErr)
			return match
		}
		return kStr
	})

	return processedLine
}
