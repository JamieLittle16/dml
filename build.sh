#!/bin/bash
set -e

# Script to build DML into a single binary

# Show usage if requested
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  (none)        Build the DML binary"
    echo "  install       Install DML system-wide (requires sudo)"
    echo "  install local Install DML to user's ~/.local directory"
    echo "  test          Run all tests"
    echo "  test-nomathjax Run tests that don't require LaTeX/ImageMagick"
    echo "  help          Show this help message"
    exit 0
fi

# Navigate to project root
cd "$(dirname "$0")"

# Process test command
if [ "$1" = "test" ]; then
    echo "Running all tests..."
    go test -v ./...
    exit $?
fi

# Process test-nomathjax command
if [ "$1" = "test-nomathjax" ]; then
    echo "Running tests without LaTeX dependencies..."
    SKIP_LATEX_TESTS=1 go test -v ./...
    exit $?
fi

# Build the binary (default action)
if [ "$1" = "" ] || [ "$1" = "build" ]; then
    echo "Building DML..."
    go build -o dml cmd/dml/main.go
    echo "Build successful! Binary created at: ./dml"
fi

# Check for install command
if [ "$1" = "install" ]; then
    echo "Installing DML..."
    
    # Build first if the binary doesn't exist
    if [ ! -f "./dml" ]; then
        echo "Building DML first..."
        go build -o dml cmd/dml/main.go
    fi
    
    # Determine install location
    if [ "$2" = "local" ]; then
        # Local installation
        mkdir -p $HOME/.local/bin
        cp dml $HOME/.local/bin/dml
        echo "DML installed to $HOME/.local/bin/dml"
        
        # Install man page
        mkdir -p $HOME/.local/share/man/man1
        cp man/dml.1 $HOME/.local/share/man/man1/dml.1
        gzip -f $HOME/.local/share/man/man1/dml.1
        echo "Man page installed to $HOME/.local/share/man/man1/dml.1.gz"
        echo "Ensure $HOME/.local/bin is in your PATH and $HOME/.local/share/man is in your MANPATH."
    else
        # System-wide installation (requires sudo)
        sudo cp dml /usr/local/bin/dml
        echo "DML installed to /usr/local/bin/dml"
        
        # Install man page
        sudo mkdir -p /usr/local/share/man/man1
        sudo cp man/dml.1 /usr/local/share/man/man1/dml.1
        sudo gzip -f /usr/local/share/man/man1/dml.1
        echo "Man page installed to /usr/local/share/man/man1/dml.1.gz"
        echo "You might need to run 'sudo mandb' for the system to recognize the new man page."
    fi
fi

# Only print "Done!" for commands that don't have their own output
if [ "$1" = "" ] || [ "$1" = "build" ] || [ "$1" = "install" ]; then
    echo "Done!"
fi