package model

import (
	"fmt"
	"os"
	"path/filepath"
)

type SourceCode struct {
	Filename string
	Language string
	Content  string
}

func (c SourceCode) Save(directory string) error {
	fmt.Fprintf(os.Stderr, "Saving file: %s in %s", c.Filename, directory)
	return os.WriteFile(filepath.Join(directory, c.Filename), []byte(c.Content), 0644)
}

func (c SourceCode) String() string {
	return c.Content
}
