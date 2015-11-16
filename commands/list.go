package commands

import (
	"log"

	"github.com/gophersaurus/govend/packages"
	"github.com/spf13/cobra"
)

var (
	listWrite  string
	listFormat string
	listAll    bool
	listVendor bool
)

func init() {
	ListCMD.Flags().StringVarP(&listWrite, "write", "w", "", "Write the results to disk")
	ListCMD.Flags().StringVarP(&listFormat, "format", "f", "txt", "Format the results with values txt, json, yml, or xml")
	ListCMD.Flags().BoolVar(&listVendor, "vendor", false, "Show all vendor dependecy packages")
	ListCMD.Flags().BoolVar(&listAll, "all", false, "Show all packages, even those in the standard library")
}

// ListCMD describes the scan command.
var ListCMD = &cobra.Command{
	Use:   "list",
	Short: "List external package dependencies.",
	Long:  "List external package dependencies in a golang project directory.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			if err := packages.ListCMD(args[0], vendorDir, listWrite, listFormat, listAll, listVendor); err != nil {
				log.Fatal(err)
			}
			return
		}

		if err := packages.ListCMD(".", vendorDir, listWrite, listFormat, listAll, listVendor); err != nil {
			log.Fatal(err)
		}
	},
}
