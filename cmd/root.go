package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vsyaco/html2md/internal"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var (
	flagArticle  bool
	flagOutput   string
	flagStdout   bool
	flagNoImages bool
	flagDomain   string
)

var rootCmd = &cobra.Command{
	Use:     "html2md [flags] <file.html>",
	Short:   "Convert HTML pages to Markdown",
	Version: fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date),
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		if _, err := os.Stat(inputPath); err != nil {
			return fmt.Errorf("file not found: %s", inputPath)
		}

		opts := internal.Options{
			Article:  flagArticle,
			NoImages: flagNoImages,
			Domain:   flagDomain,
		}

		result, err := internal.Convert(inputPath, opts)
		if err != nil {
			return err
		}

		if flagStdout {
			fmt.Print(result)
			return nil
		}

		outputPath := flagOutput
		if outputPath == "" {
			outputPath = replaceExt(inputPath, ".md")
		}

		if err := os.WriteFile(outputPath, []byte(result+"\n"), 0644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Written to %s\n", outputPath)
		return nil
	},
}

// Execute runs the root command.
func Execute(version, commit, date string) {
	Version = version
	Commit = commit
	Date = date
	rootCmd.Version = fmt.Sprintf("%s (commit %s, built %s)", version, commit, date)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&flagArticle, "article", false, "Extract main content only (readability)")
	rootCmd.Flags().StringVarP(&flagOutput, "output", "o", "", "Output .md file path")
	rootCmd.Flags().BoolVar(&flagStdout, "stdout", false, "Print to stdout instead of file")
	rootCmd.Flags().BoolVar(&flagNoImages, "no-images", false, "Remove images")
	rootCmd.Flags().StringVar(&flagDomain, "domain", "", "Base URL for resolving relative links")
}

func replaceExt(path, ext string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(dir, name+ext)
}
