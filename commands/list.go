package commands

import (
	"log"

	"github.com/gophersaurus/govend/govend"
	"github.com/spf13/cobra"
)

var (
	write  string
	all    bool
	vendor bool
)

func init() {
	ListCMD.Flags().StringVarP(&write, "write", "w", "", "Write the results to disk")
	ListCMD.Flags().StringVarP(&format, "format", "f", "txt", "Format the results with values txt, json, yml, or xml")
	ListCMD.Flags().BoolVar(&vendor, "vendor", false, "Show all vendor dependecy packages")
	ListCMD.Flags().BoolVar(&all, "all", false, "Show all packages, even those in the standard library")
}

// ListCMD describes the scan command.
var ListCMD = &cobra.Command{
	Use:   "list",
	Short: "List external package dependencies.",
	Long:  "List external package dependencies in a golang project directory.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := govend.List(args, write, format, all, vendor); err != nil {
			log.Fatal(err)
		}
		return
	},
}
