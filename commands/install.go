package commands

import (
	"log"

	"github.com/gophersaurus/govend/repo"
	"github.com/spf13/cobra"
)

// InstallCMD describes the install command.
var InstallCMD = &cobra.Command{
	Use:   "install",
	Short: "Download a project's external package dependencies.",
	Long:  "Download a project's external repository dependencies by what is provided in the vendor manifest file.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := repo.InstallCMD(vendorDir, vendorFile, verbose); err != nil {
			log.Fatal(err)
		}
	},
}
