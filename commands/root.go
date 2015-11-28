package commands

import (
	"log"

	"github.com/gophersaurus/govend/govend"
	"github.com/spf13/cobra"
)

var (
	verbose  bool
	update   bool
	commands bool
	format   string
)

const (
	rootDesc = `Govend downloads and vendors the packages named by the import
paths, along with their dependencies.`
	verboseDesc = `The -v flag prints the names of packages as they are vendored.
	`
	updateDesc = `The -u flag uses the network to update the named packages and
	their dependencies.  By default, the network is used to check out missing
	packages but does not use it to look for updates to existing packages.
	`
	commandsDesc = `The -x flag prints commands as they are executed for vendoring
	such as 'git init'.
	`
	formatDesc = `The -f flag defines the format of manifest vendor file on disk.
	By default, the file format is YAML but also supports JSON and TOML formats.
	`
)

func init() {
	RootCMD.Flags().BoolVarP(&commands, "commands", "x", false, commandsDesc)
	RootCMD.Flags().StringVarP(&format, "format", "f", "YAML", formatDesc)
	RootCMD.Flags().BoolVarP(&update, "update", "u", false, updateDesc)
	RootCMD.Flags().BoolVarP(&verbose, "verbose", "v", false, verboseDesc)
}

// RootCMD describes the root command.
var RootCMD = &cobra.Command{
	Short: "Govend vendors external packages.",
	Long:  rootDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if err := govend.Vendor(args, update, verbose, commands, format); err != nil {
			log.Fatal(err)
		}
	},
}
