#!/bin/bash
set -e

# Script to build DML into a single binary

echo "Building DML..."

# Navigate to project root
cd "$(dirname "$0")"

# Build the binary
go build -o dml cmd/dml/main.go

echo "Build successful! Binary created at: ./dml"

# Check for install command
if [ "$1" = "install" ]; then
    echo "Installing DML..."
    
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

echo "Done!"