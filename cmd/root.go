/*
Copyright © 2023 richard.wooding@spandigital.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/spandigitial/codeblocks/model"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codeblocks",
	Short: "Extract fenced code blocks from markdown",
	Long:  `Extracts fenced code blocks from markdown`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input := viper.GetString("input")
		var source []byte
		var err error
		if input == "" {
			source, err = io.ReadAll(os.Stdin)
		} else {
			source, err = os.ReadFile(input)
		}
		if err != nil {
			return err
		}

		extension := viper.GetString("extension")
		if extension == "" {
			extension = "txt"
		}
		filenamePrefix := viper.GetString("filename-prefix")
		if filenamePrefix == "" {
			filenamePrefix = "sourcecode"
		}

		outputDirectory := viper.GetString("output-directory")
		if outputDirectory == "" {
			outputDirectory, err = os.Getwd()
			if err != nil {
				return err
			}
		}

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

		l := len(codeBlocks)
		userSpecifiedExtension := viper.GetString("extension") != "" // Check if user provided --extension

		for i, codeBlock := range codeBlocks {
			sourceCode := codeBlock.ToSourceCode(func(block model.FencedCodeBlock) string {
				// Determine extension: user override > language detection > default fallback
				fileExtension := extension // Default
				if !userSpecifiedExtension {
					// Auto-detect extension from language (handles empty strings)
					fileExtension = model.LanguageToExtension(block.Language)
				}

				if l == 1 {
					return fmt.Sprintf("%s.%s", filenamePrefix, fileExtension)
				} else {
					return fmt.Sprintf("%s-%d.%s", filenamePrefix, i, fileExtension)
				}
			})
			if err := sourceCode.Save(outputDirectory); err != nil {
				return fmt.Errorf("failed to save %s: %w", sourceCode.Filename, err)
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codeblocks.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("input", "i", "", "Input (defaults to stdin)")
	if err := viper.BindPFlag("input", rootCmd.Flags().Lookup("input")); err != nil {
		log.Fatal("Unable to bind flag input", err)
	}
	rootCmd.Flags().StringP("extension", "e", "", "Extension (defaults to txt)")
	if err := viper.BindPFlag("extension", rootCmd.Flags().Lookup("extension")); err != nil {
		log.Fatal("Unable to bind flag extension", err)
	}
	rootCmd.Flags().StringP("filename-prefix", "f", "", "Filename prefix (defaults to sourcecode)")
	if err := viper.BindPFlag("filename-prefix", rootCmd.Flags().Lookup("filename-prefix")); err != nil {
		log.Fatal("Unable to bind filename-prefix", err)
	}
	rootCmd.Flags().StringP("output-directory", "o", "", "Output directory (defaults to current working directory)")
	if err := viper.BindPFlag("output-directory", rootCmd.Flags().Lookup("output-directory")); err != nil {
		log.Fatal("Unable to bind flag output-directory", err)
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".codeblocks" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".codeblocks")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
