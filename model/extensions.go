package model

import "strings"

// LanguageToExtension maps common programming language identifiers to file extensions.
// It performs case-insensitive matching and returns "txt" for unknown languages.
func LanguageToExtension(language string) string {
	// Map common language identifiers to extensions
	extensionMap := map[string]string{
		// Compiled languages
		"go":      "go",
		"golang":  "go",
		"rust":    "rs",
		"c":       "c",
		"cpp":     "cpp",
		"c++":     "cpp",
		"java":    "java",
		"kotlin":  "kt",
		"swift":   "swift",
		"csharp":  "cs",
		"c#":      "cs",
		"objc":    "m",
		"haskell": "hs",
		"scala":   "scala",

		// Scripting languages
		"python":  "py",
		"python3": "py",
		"ruby":    "rb",
		"perl":    "pl",
		"php":     "php",
		"lua":     "lua",
		"r":       "R",
		"julia":   "jl",

		// Web languages
		"javascript": "js",
		"js":         "js",
		"typescript": "ts",
		"ts":         "ts",
		"html":       "html",
		"css":        "css",
		"scss":       "scss",
		"sass":       "sass",
		"less":       "less",
		"jsx":        "jsx",
		"tsx":        "tsx",
		"vue":        "vue",
		"svelte":     "svelte",

		// Shell
		"bash":       "sh",
		"sh":         "sh",
		"shell":      "sh",
		"zsh":        "sh",
		"fish":       "fish",
		"powershell": "ps1",
		"ps1":        "ps1",

		// Data formats
		"json":       "json",
		"yaml":       "yaml",
		"yml":        "yaml",
		"toml":       "toml",
		"xml":        "xml",
		"ini":        "ini",
		"properties": "properties",

		// Markup
		"markdown": "md",
		"md":       "md",
		"tex":      "tex",
		"latex":    "tex",

		// Database
		"sql":        "sql",
		"postgres":   "sql",
		"postgresql": "sql",
		"mysql":      "sql",
		"sqlite":     "sql",
		"plsql":      "sql",
		"tsql":       "sql",

		// Other
		"dockerfile": "Dockerfile",
		"docker":     "Dockerfile",
		"makefile":   "Makefile",
		"make":       "Makefile",
		"graphql":    "graphql",
		"protobuf":   "proto",
		"proto":      "proto",
		"diff":       "diff",
		"patch":      "patch",
	}

	// Convert to lowercase for case-insensitive matching
	ext, found := extensionMap[strings.ToLower(language)]
	if found {
		return ext
	}

	// Fallback: use txt for unknown languages
	return "txt"
}
