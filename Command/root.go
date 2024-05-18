package Command

import (
	"gport/Common"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gport",
	Example: "gport share -a localhost:80",
	Version: Common.CliVersion,
	Short:   "Connect Your Local World to the Global Internet.",
	Long:    `gPort is a reverse proxy service that allows your local addresses to be accessible from the global internet. Easily and securely connect your local servers and applications to the outside world with unique global URLs provided by gPort.`,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("gPort CLI version: {{.Version}}")
}
