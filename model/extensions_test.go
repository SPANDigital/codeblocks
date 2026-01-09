package model

import "testing"

func TestLanguageToExtension(t *testing.T) {
	tests := []struct {
		language string
		expected string
	}{
		// Compiled languages
		{"go", "go"},
		{"golang", "go"},
		{"Go", "go"},      // Case insensitive
		{"GO", "go"},      // All caps
		{"rust", "rs"},
		{"c", "c"},
		{"cpp", "cpp"},
		{"c++", "cpp"},
		{"java", "java"},
		{"kotlin", "kt"},

		// Scripting languages
		{"python", "py"},
		{"python3", "py"},
		{"Python", "py"},  // Case insensitive
		{"ruby", "rb"},
		{"perl", "pl"},
		{"php", "php"},

		// Web languages
		{"javascript", "js"},
		{"js", "js"},
		{"JavaScript", "js"},  // Case insensitive
		{"typescript", "ts"},
		{"ts", "ts"},
		{"html", "html"},
		{"css", "css"},
		{"jsx", "jsx"},
		{"tsx", "tsx"},

		// Shell
		{"bash", "sh"},
		{"sh", "sh"},
		{"shell", "sh"},
		{"zsh", "sh"},
		{"Bash", "sh"},    // Case insensitive
		{"powershell", "ps1"},

		// Data formats
		{"json", "json"},
		{"yaml", "yaml"},
		{"yml", "yaml"},
		{"toml", "toml"},
		{"xml", "xml"},

		// Markup
		{"markdown", "md"},
		{"md", "md"},

		// Database
		{"sql", "sql"},
		{"postgres", "sql"},
		{"mysql", "sql"},
		{"postgresql", "sql"},

		// Other
		{"dockerfile", "Dockerfile"},
		{"docker", "Dockerfile"},
		{"Dockerfile", "Dockerfile"},  // Case insensitive
		{"makefile", "Makefile"},
		{"make", "Makefile"},

		// Fallback cases
		{"unknown", "txt"},          // Unknown language
		{"foobar", "txt"},           // Random string
		{"", "txt"},                 // Empty string
		{"NOT_A_REAL_LANGUAGE", "txt"},  // All caps unknown
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			result := LanguageToExtension(tt.language)
			if result != tt.expected {
				t.Errorf("LanguageToExtension(%q) = %q, want %q",
					tt.language, result, tt.expected)
			}
		})
	}
}

// TestLanguageToExtensionConsistency ensures that common language aliases map to the same extension
func TestLanguageToExtensionConsistency(t *testing.T) {
	aliasGroups := [][]string{
		{"go", "golang"},
		{"javascript", "js"},
		{"typescript", "ts"},
		{"python", "python3"},
		{"bash", "sh", "shell"},
		{"markdown", "md"},
		{"yaml", "yml"},
		{"cpp", "c++"},
	}

	for _, group := range aliasGroups {
		var expectedExt string
		for i, lang := range group {
			ext := LanguageToExtension(lang)
			if i == 0 {
				expectedExt = ext
			} else if ext != expectedExt {
				t.Errorf("Inconsistent mapping: %q -> %q, but %q -> %q",
					group[0], expectedExt, lang, ext)
			}
		}
	}
}
