// Package cmd define subcommands provided by the gup command
package cmd

import (
	"os"

	"github.com/nao1215/gup/internal/assets"
	"github.com/nao1215/gup/internal/completion"
	"github.com/nao1215/gup/internal/print"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gup",
	Short: `gup command update binaries installed by 'go install'.
If you update all binaries, just run '$ gup update'`,
}

// OsExit is wrapper for  os.Exit(). It's for unit test.
var OsExit = os.Exit

// Execute run gup process.
func Execute() {
	assets.DeployIconIfNeeded()
	completion.DeployShellCompletionFileIfNeeded(rootCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		print.Err(err)
	}
}
