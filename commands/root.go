package commands

import (
	"log"

	"github.com/gophersaurus/govend/repo"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	recursive  bool
	vendorDir  string
	vendorFile string
)

func init() {
	RootCMD.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output")
	RootCMD.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Perform the command recurively if possible.")
	RootCMD.PersistentFlags().StringVar(&vendorDir, "vendorDir", "vendor", "Define the vendor directory")
	RootCMD.PersistentFlags().StringVar(&vendorFile, "vendorFile", "vendor.yml", "Define the vendor manifest file")
}

// RootCMD describes the root command.
var RootCMD = &cobra.Command{
	Short: "Vendor a project's external package dependencies.",
	Long:  "Vendor a project's external package dependencies.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := repo.VendCMD(vendorDir, vendorFile, verbose, recursive); err != nil {
			log.Fatal(err)
		}
	},
}
