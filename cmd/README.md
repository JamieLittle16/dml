# Command Package

This directory contains the command-line application entry points for DML.

## Structure

- `dml/` - The main DML application entry point
  - `main.go` - Contains the main application logic, command-line flag parsing, and orchestration of the processing pipeline

The `main.go` file is responsible for:
1. Parsing command-line flags and options
2. Setting up the processing environment
3. Coordinating between different processing modes (streaming vs. full document)
4. Managing the input/output streams
5. Error handling and reporting

This separation of concerns allows the command-line interface to evolve independently from the core processing logic contained in the internal packages.