# codeblocks

A command-line tool to extract fenced code blocks from markdown files and save them as individual source files.

## Features

- **Automatic extension detection** - Files get appropriate extensions based on code language (Go → `.go`, Python → `.py`, etc.)
- Extract code blocks from markdown files or stdin
- Preserve language information from fenced code blocks
- Customize output filenames and extensions
- Support for multiple code blocks in a single markdown file
- Configurable via CLI flags, config file, or environment variables

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap spandigital/tap
brew install codeblocks
```

### Go Install

```bash
go install github.com/spandigitial/codeblocks@latest
```

### From Source

```bash
git clone https://github.com/SPANDigital/codeblocks.git
cd codeblocks
go build
```

## Usage

### Basic Examples

Extract code blocks from stdin:

```bash
cat README.md | codeblocks
```

Extract code blocks from a file:

```bash
codeblocks -i documentation.md
```

Specify output directory and file extension:

```bash
codeblocks -i tutorial.md -o ./examples -e go -f example
```

### Real-World Use Cases

**Extract code examples from documentation:**
```bash
# Extract all code blocks from API documentation
codeblocks -i api-docs.md -o ./code-samples -f api-example
```

**Process markdown from a URL via pipe:**
```bash
# Download and extract code blocks
curl -s https://raw.githubusercontent.com/user/repo/main/README.md | codeblocks -o ./extracted
```

**Extract code for testing:**
```bash
# Extract test examples from documentation
codeblocks -i tests.md -e test.go -f integration -o ./test-cases
```

**Batch process multiple files:**
```bash
# Extract code from all markdown files
for file in docs/*.md; do
  codeblocks -i "$file" -o ./code-samples -f "$(basename "$file" .md)"
done
```

## Language-Based File Extensions

By default, `codeblocks` automatically detects the programming language from fenced code blocks and uses the appropriate file extension. This means your extracted code files will have the correct extension for their language, making them immediately usable.

### Automatic Detection

```bash
# Input markdown with different languages
$ cat example.md
```go
package main

func main() {
    println("Hello from Go!")
}
```

```python
def greet():
    print("Hello from Python!")
```

```javascript
function greet() {
    console.log("Hello from JavaScript!");
}
```

# Extract with auto-detected extensions
$ codeblocks -i example.md
Saving file: sourcecode-0.go in /current/directory
Saving file: sourcecode-1.py in /current/directory
Saving file: sourcecode-2.js in /current/directory
```

### Supported Languages

The tool automatically recognizes 40+ programming languages and data formats:

- **Compiled languages:** Go (`.go`), Rust (`.rs`), C (`.c`), C++ (`.cpp`), Java (`.java`), Kotlin (`.kt`), Swift (`.swift`)
- **Scripting languages:** Python (`.py`), Ruby (`.rb`), Perl (`.pl`), PHP (`.php`), Lua (`.lua`)
- **Web technologies:** JavaScript (`.js`), TypeScript (`.ts`), HTML (`.html`), CSS (`.css`), JSX (`.jsx`), TSX (`.tsx`)
- **Shell scripts:** Bash/Shell (`.sh`), Fish (`.fish`), PowerShell (`.ps1`)
- **Data formats:** JSON (`.json`), YAML (`.yaml`), TOML (`.toml`), XML (`.xml`)
- **Markup:** Markdown (`.md`), LaTeX (`.tex`)
- **Database:** SQL (`.sql`)
- **Other:** Dockerfile, Makefile, and more...

### Override Auto-Detection

If you need all files to have the same extension, use the `--extension` flag to override auto-detection:

```bash
# Force all code blocks to use .txt extension
$ codeblocks -i example.md --extension txt
Saving file: sourcecode-0.txt
Saving file: sourcecode-1.txt
Saving file: sourcecode-2.txt
```

This is useful when:
- You want uniform extensions regardless of language
- You're extracting code snippets for documentation
- You need compatibility with systems that expect specific extensions

### Unknown Languages

Code blocks with unknown or missing language identifiers automatically fallback to `.txt`:

```bash
# Markdown with unknown language
$ cat example.md
```unknownlang
some code in an unrecognized language
```

# Output uses .txt fallback
$ codeblocks -i example.md
Saving file: sourcecode.txt
```

## Command-Line Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--input` | `-i` | Input markdown file | stdin |
| `--extension` | `-e` | File extension for output files (overrides auto-detection) | Auto-detected from language |
| `--filename-prefix` | `-f` | Prefix for output filenames | `sourcecode` |
| `--output-directory` | `-o` | Output directory | Current directory |
| `--config` | | Config file path | `$HOME/.codeblocks.yaml` |
| `--help` | `-h` | Show help information | |

## Configuration

You can configure `codeblocks` using:

1. **Command-line flags** (highest priority)
2. **Environment variables** (prefix with `CODEBLOCKS_`, e.g., `CODEBLOCKS_EXTENSION=go`)
3. **Config file** at `$HOME/.codeblocks.yaml` (lowest priority)

### Example Config File

Create `$HOME/.codeblocks.yaml`:

```yaml
extension: go
filename-prefix: example
output-directory: ./code-samples
```

## How It Works

`codeblocks` parses markdown using [goldmark](https://github.com/yuin/goldmark), walks the AST to find fenced code blocks, and extracts them with their language information. Each code block is saved as a separate file with the specified prefix and extension.

**Example Input:**

````markdown
# My Tutorial

Here's a Go example:

```go
package main

func main() {
    println("Hello, World!")
}
```

And a Python example:

```python
def hello():
    print("Hello, World!")
```
````

**Command:**
```bash
codeblocks -i tutorial.md -f example
```

**Output (with automatic extension detection):**
- `example-0.go` (contains the Go code)
- `example-1.py` (contains the Python code)

## Development

### Prerequisites

- Go 1.25 or higher

### Building

```bash
go build -v ./...
```

### Running Tests

```bash
go test -v ./...
```

### Running Locally

```bash
echo '```go
package main
func main() { println("test") }
```' | go run main.go -e go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Links

- [GitHub Repository](https://github.com/SPANDigital/codeblocks)
- [Issue Tracker](https://github.com/SPANDigital/codeblocks/issues)
- [Releases](https://github.com/SPANDigital/codeblocks/releases)
