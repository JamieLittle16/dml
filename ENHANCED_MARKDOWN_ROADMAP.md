# Enhanced Markdown Support Roadmap

This document tracks DML's progress toward feature parity with [Glow](https://github.com/charmbracelet/glow) and full CommonMark/GitHub Flavored Markdown (GFM) specification compliance.

## Current Status

### ✅ Implemented Features

**Basic Markdown:**
- [x] **Bold text** (`**bold**`, `__bold__`)
- [x] *Italic text* (`*italic*`, `_italic_`)
- [x] `Inline code` (`` `code` ``)
- [x] Code blocks (`` ``` ``)
- [x] Headings (`# H1` through `###### H6`)
- [x] Paragraphs and line breaks
- [x] Horizontal rules (`---`)

**LaTeX Integration:**
- [x] Inline math (`$formula$`)
- [x] Display math (`$$formula$$`)
- [x] Mixed Markdown + LaTeX document rendering

### 🚧 In Progress / Scaffolded

**Tables:**
- [ ] Basic table syntax (`| col1 | col2 |`)
- [ ] Table headers with alignment
- [ ] Terminal ANSI table rendering
- [ ] LaTeX table generation

**Lists:**
- [ ] Unordered lists (`-`, `*`, `+`)
- [ ] Ordered lists (`1.`, `2.`, etc.)
- [ ] Nested lists
- [ ] List item formatting

**Block Elements:**
- [ ] Blockquotes (`> quote`)
- [ ] Nested blockquotes

**Links and References:**
- [ ] Inline links (`[text](url)`)
- [ ] Reference links (`[text][ref]`)
- [ ] Autolinks (`<url>`)
- [ ] Clickable terminal links (where supported)

**Images:**
- [ ] Image syntax (`![alt](src)`)
- [ ] Alt text rendering as placeholder
- [ ] Future: Inline image display via terminal protocols

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