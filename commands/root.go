package commands

import (
	"log"
	"path/filepath"

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
	vfile := filepath.Join("vendor", "vendor.yml")
	RootCMD.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output to os.Stdout.")
	RootCMD.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Execute the command recurively.")
	RootCMD.PersistentFlags().StringVar(&vendorDir, "vendorDir", "vendor", "Define the vendor directory location on disk.")
	RootCMD.PersistentFlags().StringVar(&vendorFile, "vendorFile", vfile, "Define the vendor manifest file location on disk.")
}

// RootCMD describes the root command.
var RootCMD = &cobra.Command{
	Short: "Vendor external packages.",
	Long:  "Vendor a Go project's external dependent packages.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := repo.VendCMD(verbose); err != nil {
			log.Fatal(err)
		}
	},
}
