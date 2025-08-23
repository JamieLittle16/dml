# Enhanced Markdown Support Roadmap

This document tracks DML's progress toward feature parity with [Glow](https://github.com/charmbracelet/glow) and full CommonMark/GitHub Flavored Markdown (GFM) specification compliance.

## Current Status

### ✅ Implemented Features

**Basic Markdown:**
- [x] **Bold text** (`**bold**`, `__bold__`) with ANSI bold formatting
- [x] *Italic text* (`*italic*`, `_italic_`) with ANSI italic formatting
- [x] `Inline code` (`` `code` ``) with ANSI reverse video
- [x] Code blocks (`` ``` ``) with ANSI reverse video
- [x] Headings (`# H1` through `###### H6`) with ANSI bold + underline
- [x] Paragraphs and line breaks
- [x] Horizontal rules (`---`) with styled terminal display

**LaTeX Integration:**
- [x] Inline math (`$formula$`)
- [x] Display math (`$$formula$$`)
- [x] Mixed Markdown + LaTeX document rendering

**Enhanced Parser Support:**
- [x] GFM (GitHub Flavored Markdown) parser extensions enabled
- [x] CommonMark compliance extensions
- [x] Tables, Strikethrough, FencedCode, Autolink extensions active

### 🚧 In Progress / Scaffolded

**Tables:**
- [x] Parser extension enabled for table recognition
- [x] Basic table AST node handling (scaffolded)
- [ ] Complete terminal table rendering with borders and alignment
- [ ] LaTeX table generation with proper formatting

**Lists:**
- [x] Basic unordered list bullets (`-`, `*`, `+`)
- [x] Basic ordered list numbering (`1.`, `2.`, etc.)
- [x] AST node handling for lists and list items
- [ ] Proper nested list indentation
- [ ] Advanced list item formatting

**Block Elements:**
- [x] Basic blockquote formatting (`> quote`) with quote bar
- [x] Horizontal rule styling
- [ ] Nested blockquotes
- [ ] Enhanced blockquote styling

**Links and References:**
- [x] Basic link text extraction (`[text](url)`)
- [x] Link styling with underline and color
- [ ] Clickable terminal links (OSC 8 escape codes)
- [ ] Reference links (`[text][ref]`)
- [ ] Autolinks (`<url>`)

**Images:**
- [x] Image placeholder rendering (`![alt](src)`)
- [x] Alt text display
- [ ] Enhanced placeholder formatting
- [ ] Future: Inline image display via terminal protocols

**Text Formatting:**
- [x] ~~Strikethrough~~ AST node support (`~~text~~`)
- [x] Scaffolded ANSI strikethrough formatting
- [ ] Verify strikethrough rendering functionality
- [ ] Additional text decoration support

### 📋 Planned Features (GFM Extensions)

**Task Lists:**
- [ ] Task list items (`- [ ]` unchecked, `- [x]` checked)
- [ ] Interactive checkbox rendering

**Text Formatting:**
- [ ] ~~Strikethrough~~ (`~~text~~`)
- [ ] Underline support

**Advanced Features:**
- [ ] Emoji support (`:emoji:` syntax)
- [ ] Syntax highlighting for code blocks
- [ ] Language detection and highlighting
- [ ] Table of contents generation
- [ ] Footnotes
- [ ] Definition lists

**HTML Compatibility:**
- [ ] Basic HTML tag support
- [ ] HTML entity rendering

## Target Compatibility

### Glow Feature Parity
- [ ] All Glow-supported Markdown features
- [ ] Similar terminal rendering quality
- [ ] Theme and styling compatibility

### CommonMark Compliance
- [ ] Full CommonMark 0.30 specification
- [ ] Proper edge case handling
- [ ] Unicode support

### GitHub Flavored Markdown (GFM)
- [ ] Tables
- [ ] Task lists
- [ ] Strikethrough
- [ ] Autolinks
- [ ] Fenced code blocks with language info

## Implementation Strategy

### Phase 1: Core Structure (Current)
- ✅ Basic parser integration
- ✅ AST traversal framework
- 🚧 Feature scaffolding and stubs

### Phase 2: Essential Features
- Tables (high priority for documentation)
- Lists (fundamental markdown feature)
- Blockquotes (common in documentation)

### Phase 3: Enhanced Features
- Links with terminal clickability
- Image placeholder rendering
- Syntax highlighting foundation

### Phase 4: GFM Extensions
- Task lists
- Strikethrough
- Emoji support
- Advanced table features

### Phase 5: Advanced Features
- Interactive elements
- Performance optimization
- Extended theme support

## Technical Notes

### Parser Extensions
Using `github.com/gomarkdown/markdown` with extensions:
- CommonMark compliance
- GFM table support
- Task list extension
- Strikethrough extension
- Autolink extension

### Rendering Targets
- **Terminal (ANSI)**: Rich text with colors, formatting, and (where supported) clickable links
- **LaTeX**: High-quality document generation for full-document rendering mode

### Dependencies
- Core: `github.com/gomarkdown/markdown`
- Terminal: Kitty graphics protocol support
- LaTeX: Standard LaTeX packages (tables, listings, etc.)

## Contributing

When implementing new features:
1. Add test cases covering the new syntax
2. Implement both terminal (ANSI) and LaTeX rendering
3. Update this roadmap with progress
4. Ensure backward compatibility
5. Add documentation and examples

## References

- [CommonMark Specification](https://commonmark.org/spec/)
- [GitHub Flavored Markdown Spec](https://github.github.com/gfm/)
- [Glow Project](https://github.com/charmbracelet/glow)
- [gomarkdown Documentation](https://github.com/gomarkdown/markdown)