#!/bin/bash

# Enhanced Markdown Support Demo for DML
# This script demonstrates the scaffolded enhanced markdown features

echo "=== DML Enhanced Markdown Support Demo ==="
echo

echo "1. Basic formatting (working):"
echo "This is **bold** and *italic* text" | ./dml
echo

echo "2. Lists (scaffolded):"
echo -e "- Item 1\n- Item 2\n- Item 3" | ./dml
echo

echo "3. Blockquotes (scaffolded):"
echo "> This is a blockquote" | ./dml
echo

echo "4. Links (scaffolded):"
echo "[Example Link](https://example.com)" | ./dml
echo

echo "5. Images (scaffolded):"
echo "![Alt Text](image.png)" | ./dml
echo

echo "6. Code blocks (enhanced):"
echo -e "\`\`\`javascript\nconsole.log('Hello, World!');\n\`\`\`" | ./dml
echo

echo "7. Horizontal rules (enhanced):"
echo "---" | ./dml
echo

echo "8. Tables (scaffolded):"
echo -e "| Name | Age |\n|------|-----|\n| John | 30  |" | ./dml
echo

echo "=== Demo Complete ==="
echo "See ENHANCED_MARKDOWN_ROADMAP.md for detailed feature status"