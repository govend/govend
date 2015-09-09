package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

func init() { version = "govend version 0.5.0" }

// VersionCMD describes the scan command.
var VersionCMD = &cobra.Command{
	Use:   "version",
	Short: "Installed version of govend.",
	Long:  "Print the current installed version of govend.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
