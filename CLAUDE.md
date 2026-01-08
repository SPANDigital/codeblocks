# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`codeblocks` is a Go CLI tool that extracts fenced code blocks from markdown files. It uses the goldmark library for markdown parsing and cobra for the CLI framework.

**Core functionality**: Parses markdown input (from stdin or file) and extracts all fenced code blocks, saving them as individual files.

## Architecture

### Entry Point
- `main.go` - Minimal entry point that delegates to cmd package

### Command Layer (`cmd/`)
- `cmd/root.go` - Cobra command setup and main application logic
  - Parses markdown using goldmark's AST walker
  - Extracts `FencedCodeBlock` objects by walking the AST
  - Converts blocks to `SourceCode` and saves to disk
  - Configuration via viper (supports flags, config file at `$HOME/.codeblocks.yaml`, and env vars)

### Model Layer (`model/`)
- `model/fencedcodeblock.go` - Represents extracted markdown code blocks (language + content)
- `model/sourcecode.go` - Represents output files (filename + language + content) with Save method

**Key Pattern**: `FencedCodeBlock` → `ToSourceCode()` → `SourceCode.Save()` pipeline

## Development Commands

### Build
```bash
go build -v ./...
```

### Run Tests
```bash
go test -v ./...
```

### Run Locally
```bash
# From stdin
echo '```go\nfmt.Println("hello")\n```' | go run main.go

# From file
go run main.go -i input.md -e go -f output -o ./out
```

### Build for Release
```bash
# Requires goreleaser installed
goreleaser release --snapshot --clean
```

## CLI Flags

- `-i, --input`: Input file (defaults to stdin)
- `-e, --extension`: File extension for output files (default: "txt")
- `-f, --filename-prefix`: Prefix for output filenames (default: "sourcecode")
- `-o, --output-directory`: Output directory (default: current working directory)
- `--config`: Config file path (default: `$HOME/.codeblocks.yaml`)

## Release Process

Releases are automated via GitHub Actions:
- Push a tag to trigger `.github/workflows/release.yml`
- GoReleaser builds binaries for Linux, Windows, Darwin
- Creates archives and checksums
- Updates Homebrew tap at `spandigital/homebrew-tap`

## Dependencies

- `github.com/yuin/goldmark` - Markdown parsing (AST-based)
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management

## Known Issues

- [Issue #5](https://github.com/SPANDigital/codeblocks/issues/5) - Typo in `output-directory` flag and missing viper binding (cmd/root.go:156)
