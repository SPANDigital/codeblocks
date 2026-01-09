package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spandigitial/codeblocks/model"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Test helper functions

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "codeblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

func cleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup test dir %s: %v", dir, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(content)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func extractCodeBlocks(t *testing.T, markdown string) []model.FencedCodeBlock {
	t.Helper()
	source := []byte(markdown)
	node := goldmark.DefaultParser().Parse(text.NewReader(source))
	var codeBlocks []model.FencedCodeBlock

	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node.Kind() == ast.KindFencedCodeBlock {
			var language string
			var content string

			fcb := node.(*ast.FencedCodeBlock)
			if !entering && fcb.Info != nil {
				segment := fcb.Info.Segment
				language = string(source[segment.Start:segment.Stop])
				var sb strings.Builder
				lines := fcb.BaseBlock.Lines()
				l := lines.Len()
				for i := 0; i < l; i++ {
					line := lines.At(i)
					sb.Write(line.Value(source))
				}
				content = sb.String()
				if language != "" && content != "" {
					codeBlocks = append(codeBlocks, model.FencedCodeBlock{
						Language: language,
						Content:  content,
					})
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return codeBlocks
}

// Test cases

func TestSingleCodeBlock(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	markdown := `# Test
` + "```go\n" + `package main
func main() { println("Hello") }
` + "```\n"

	codeBlocks := extractCodeBlocks(t, markdown)
	if len(codeBlocks) != 1 {
		t.Fatalf("Expected 1 code block, got %d", len(codeBlocks))
	}

	// Test filename generation logic
	filenamePrefix := "sourcecode"
	extension := "txt"
	l := len(codeBlocks)

	for i, codeBlock := range codeBlocks {
		var expectedFilename string
		if l == 1 {
			expectedFilename = filenamePrefix + "." + extension
		} else {
			expectedFilename = filenamePrefix + "-" + string(rune('0'+i)) + "." + extension
		}

		sourceCode := codeBlock.ToSourceCode(func(block model.FencedCodeBlock) string {
			if l == 1 {
				return filenamePrefix + "." + extension
			}
			return filenamePrefix + "-" + string(rune('0'+i)) + "." + extension
		})

		if sourceCode.Filename != expectedFilename {
			t.Errorf("Expected filename %s, got %s", expectedFilename, sourceCode.Filename)
		}

		err := sourceCode.Save(testDir)
		if err != nil {
			t.Fatalf("Failed to save source code: %v", err)
		}
	}

	// Verify file exists
	expectedPath := filepath.Join(testDir, "sourcecode.txt")
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s does not exist", expectedPath)
	}

	// Verify content
	content := readFile(t, expectedPath)
	if !strings.Contains(content, "func main()") {
		t.Errorf("File content doesn't match expected code")
	}
}

func TestMultipleCodeBlocks(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	markdown := `# Test
` + "```go\n" + `package main
func main() { println("Go") }
` + "```\n" + `

Some text

` + "```python\n" + `def hello():
    print("Python")
` + "```\n" + `

More text

` + "```javascript\n" + `console.log("JavaScript");
` + "```\n"

	codeBlocks := extractCodeBlocks(t, markdown)
	if len(codeBlocks) != 3 {
		t.Fatalf("Expected 3 code blocks, got %d", len(codeBlocks))
	}

	// Test filename generation logic for multiple blocks
	filenamePrefix := "sourcecode"
	extension := "txt"
	l := len(codeBlocks)

	expectedFiles := []string{
		"sourcecode-0.txt",
		"sourcecode-1.txt",
		"sourcecode-2.txt",
	}

	for i, codeBlock := range codeBlocks {
		sourceCode := codeBlock.ToSourceCode(func(block model.FencedCodeBlock) string {
			if l == 1 {
				return filenamePrefix + "." + extension
			}
			return filenamePrefix + "-" + string(rune('0'+i)) + "." + extension
		})

		if sourceCode.Filename != expectedFiles[i] {
			t.Errorf("Block %d: Expected filename %s, got %s", i, expectedFiles[i], sourceCode.Filename)
		}

		err := sourceCode.Save(testDir)
		if err != nil {
			t.Fatalf("Failed to save source code block %d: %v", i, err)
		}
	}

	// Verify all files exist with correct content
	expectedContents := []string{"func main()", "print(\"Python\")", "console.log"}

	for i, filename := range expectedFiles {
		filePath := filepath.Join(testDir, filename)
		if !fileExists(filePath) {
			t.Errorf("Expected file %s does not exist", filePath)
			continue
		}

		content := readFile(t, filePath)
		if !strings.Contains(content, expectedContents[i]) {
			t.Errorf("File %s doesn't contain expected content: %s", filename, expectedContents[i])
		}
	}
}

func TestNoCodeBlocks(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	markdown := `# Test Document

This is just plain text with no code blocks.

Some more text here.
`

	codeBlocks := extractCodeBlocks(t, markdown)
	if len(codeBlocks) != 0 {
		t.Errorf("Expected 0 code blocks, got %d", len(codeBlocks))
	}

	// Verify no files were created
	entries, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected no files in directory, found %d", len(entries))
	}
}

func TestDifferentLanguages(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	languages := []string{"go", "python", "javascript", "bash", "rust"}
	var markdownBuilder strings.Builder

	for _, lang := range languages {
		markdownBuilder.WriteString("```" + lang + "\n")
		markdownBuilder.WriteString("code for " + lang + "\n")
		markdownBuilder.WriteString("```\n\n")
	}

	codeBlocks := extractCodeBlocks(t, markdownBuilder.String())
	if len(codeBlocks) != len(languages) {
		t.Fatalf("Expected %d code blocks, got %d", len(languages), len(codeBlocks))
	}

	for i, codeBlock := range codeBlocks {
		if codeBlock.Language != languages[i] {
			t.Errorf("Block %d: Expected language %s, got %s", i, languages[i], codeBlock.Language)
		}
	}
}

func TestCustomExtension(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	markdown := "```go\npackage main\n```\n"
	codeBlocks := extractCodeBlocks(t, markdown)

	extension := "go"
	filename := "test." + extension

	sourceCode := codeBlocks[0].ToSourceCode(func(block model.FencedCodeBlock) string {
		return filename
	})

	err := sourceCode.Save(testDir)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	expectedPath := filepath.Join(testDir, filename)
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s does not exist", expectedPath)
	}
}

func TestCustomPrefix(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	markdown := "```python\nprint('test')\n```\n"
	codeBlocks := extractCodeBlocks(t, markdown)

	prefix := "mycode"
	extension := "py"
	filename := prefix + "." + extension

	sourceCode := codeBlocks[0].ToSourceCode(func(block model.FencedCodeBlock) string {
		return filename
	})

	err := sourceCode.Save(testDir)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	expectedPath := filepath.Join(testDir, filename)
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s does not exist", expectedPath)
	}
}

func TestOutputDirectory(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	// Create a subdirectory
	subDir := filepath.Join(testDir, "output")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	markdown := "```bash\necho test\n```\n"
	codeBlocks := extractCodeBlocks(t, markdown)

	filename := "script.sh"
	sourceCode := codeBlocks[0].ToSourceCode(func(block model.FencedCodeBlock) string {
		return filename
	})

	err = sourceCode.Save(subDir)
	if err != nil {
		t.Fatalf("Failed to save to subdirectory: %v", err)
	}

	expectedPath := filepath.Join(subDir, filename)
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s does not exist", expectedPath)
	}
}

func TestSpecialCharactersInCode(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	// Use backticks to avoid escaping in the Go string
	markdown := "```python\n" +
		"print(\"Testing 'quotes' and double quotes\")\n" +
		"print(f\"Testing {interpolation}\")\n" +
		"```\n"

	codeBlocks := extractCodeBlocks(t, markdown)
	if len(codeBlocks) != 1 {
		t.Fatalf("Expected 1 code block, got %d", len(codeBlocks))
	}

	filename := "special.py"
	sourceCode := codeBlocks[0].ToSourceCode(func(block model.FencedCodeBlock) string {
		return filename
	})

	err := sourceCode.Save(testDir)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	expectedPath := filepath.Join(testDir, filename)
	content := readFile(t, expectedPath)

	// Verify special characters are preserved
	if !strings.Contains(content, "'quotes'") {
		t.Error("Single quotes not preserved")
	}
	if !strings.Contains(content, "double quotes") {
		t.Error("Double quotes not preserved")
	}
	if !strings.Contains(content, "{interpolation}") {
		t.Error("Curly braces not preserved")
	}
	if !strings.Contains(content, "print(") {
		t.Error("Parentheses not preserved")
	}
}

// Reset viper for isolated tests
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Cleanup
	viper.Reset()

	os.Exit(code)
}
