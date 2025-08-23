# Markdown Feature Roadmap

This document outlines the current status of Markdown features in DML and planned enhancements.

## Implemented Features ✅

### Basic Text Formatting
- **Bold text** (`**text**` or `__text__`) - ANSI and LaTeX ✅
- *Italic text* (`*text*` or `_text*`) - ANSI and LaTeX ✅
- ~~Strikethrough~~ (`~~text~~`) - ANSI and LaTeX ✅
- `Inline code` - ANSI and LaTeX ✅
- Code blocks - ANSI and LaTeX ✅

### Structural Elements
- Headings (`# H1`, `## H2`, etc.) - ANSI and LaTeX ✅
- Paragraphs - ANSI and LaTeX ✅
- Line breaks (soft and hard) - LaTeX ✅

### Lists
- Ordered lists (`1. item`) - ANSI and LaTeX ✅
- Unordered lists (`- item` or `* item`) - ANSI and LaTeX ✅
- Proper numbering and bullet characters ✅

### Blockquotes
- Basic blockquotes (`> text`) - ANSI and LaTeX ✅
- Unicode box-drawing prefix for ANSI display ✅
- LaTeX `quote` environment ✅

### Links and Images
- Links (`[text](url)`) - ANSI underline + URL display, LaTeX footnotes ✅
- Images (`![alt](url)`) - ANSI alt text + URL display, LaTeX alt text ✅

### Tables
- Basic table structure with headers and data rows ✅
- ASCII rendering for ANSI display ✅
- LaTeX `tabular` environment ✅

### Horizontal Rules
- Horizontal rules (`---`) - Unicode line characters for ANSI, `\hrulefill` for LaTeX ✅

### Math (Pre-existing)
- Inline math (`$formula$`) ✅
- Display math (`$$formula$$`) ✅

## Planned Features 📋

### List Enhancements
- [ ] Nested lists with proper indentation
- [ ] Task lists (`- [ ] unchecked`, `- [x] checked`)
- [ ] Task list checkbox rendering with ANSI glyphs (☐ ☑ ☒)
- [ ] Definition lists

### Advanced Text Features
- [ ] Smart typography (curly quotes, em/en dashes)
- [ ] Superscript and subscript (if supported by parser)
- [ ] Highlights/marks (`==text==`)

### Code and Syntax
- [ ] Syntax highlighting for code blocks (using Chroma library)
- [ ] Language-specific code block rendering
- [ ] Line numbers for code blocks

### Theming and Styling
- [ ] Theme abstraction layer for ANSI colors
- [ ] Configurable color schemes
- [ ] Custom ANSI formatting options
- [ ] Configurable link display style (inline URL vs footnote)

### Tables Enhancements
- [ ] Column alignment support (left, right, center)
- [ ] Advanced table formatting for ANSI display
- [ ] Table cell content formatting (bold, italic within cells)
- [ ] Column width optimization

### Links and References
- [ ] Reference-style links (`[text][ref]`)
- [ ] Footnotes with proper numbering
- [ ] Automatic link detection improvements

### Document Structure
- [ ] Table of contents generation
- [ ] Frontmatter handling (YAML, TOML)
- [ ] Document metadata support

### Images and Media
- [ ] Advanced image handling in LaTeX (using `\includegraphics`)
- [ ] Image size and positioning options
- [ ] Support for local image files
- [ ] Alt text improvements for accessibility

### Extensions
- [ ] Emoji conversion (`:smile:` → 😄)
- [ ] Mathematical expressions in tables
- [ ] Custom container blocks
- [ ] Admonitions/callouts

## Implementation Notes

### Current Approach
- Uses `gomarkdown/markdown` parser with extensions:
  - `CommonExtensions` (includes most basic features)
  - `AutoHeadingIDs` for automatic heading identification
  - `Strikethrough` for `~~text~~` support
  - `Tables` for table parsing
  - `Autolink` for automatic URL detection
- ANSI rendering uses standard terminal escape codes
- LaTeX generation uses appropriate environments and commands
- MathJax extension excluded to avoid conflicts with existing math processing

### Design Principles
- Maintain compatibility with existing math rendering pipeline
- Use Unicode characters where possible for better visual appeal
- Provide fallbacks for terminals with limited Unicode support
- Keep LaTeX output clean and compilable
- Minimize dependencies and complexity

### Testing Strategy
- Unit tests for individual AST node rendering
- Integration tests for complete document processing
- Manual testing with various terminal environments
- LaTeX compilation testing for full document mode

## Comparison with Glow

### Features DML Now Has (that Glow has)
- Basic Markdown formatting ✅
- Lists ✅
- Blockquotes ✅
- Tables ✅
- Links ✅
- Code blocks ✅
- Horizontal rules ✅

### Features Glow Has (that DML needs)
- [ ] Syntax highlighting
- [ ] Advanced theming
- [ ] Task list checkboxes
- [ ] Better table formatting
- [ ] Emoji support
- [ ] Style customization

### Unique DML Features
- LaTeX math rendering ✅
- Full document LaTeX generation ✅
- Kitty terminal graphics protocol ✅
- Mathematical expression support ✅

This roadmap will be updated as features are implemented and new requirements are identified.