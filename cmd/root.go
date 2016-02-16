package cmd

import (
	"fmt"
	"log"

	"github.com/govend/govend/govend"
	"github.com/spf13/cobra"
)

var (
	version  bool
	verbose  bool
	tree     bool
	update   bool
	results  bool
	commands bool
	lock     bool
	scan     bool
	all      bool
	format   string
)

const (
	rootDesc = `Govend downloads and vendors the packages named by the import
paths, along with their dependencies.`
	verboseDesc = `The -v flag prints the names of packages as they are vendored.
	`
	versionDesc = `The --version flag prints the current version.
	`
	treeDesc = `The -t flag works with the -v flag to print the names of packages
	as an indented tree to visualize the dependency tree.
	`
	resultsDesc = `The -r flag works with the -v flag to print a summary of the
	number of packages scanned, packages skipped, and repositories downloaded.
	`
	updateDesc = `The -u flag uses the network to update the named packages and
	their dependencies.  By default, the network is used to check out missing
	packages but does not use it to look for updates to existing packages.
	`
	commandsDesc = `The -x flag prints commands as they are executed for vendoring
	such as 'git init'.
	`
	lockDesc = `The -l flag writes a manifest vendor file on disk to lock in the
	versions of vendored dependencies.  This only needs to be done once.
	`
	scanDesc = `The -s flag scans the current or provided directory for external
	packages.
	`
	allDesc = `The -a flag works with the -s flag to show all packages, not just
	external packages.
	`
	formatDesc = `The -f flag works with the -m flag and -s flag to define the
	format when writing files to disk.  By default, the file format is YAML but
	also supports JSON and TOML formats.
	`
)

func init() {
	RootCMD.Flags().BoolVarP(&commands, "commands", "x", false, commandsDesc)
	RootCMD.Flags().StringVarP(&format, "format", "f", "YAML", formatDesc)
	RootCMD.Flags().BoolVarP(&update, "update", "u", false, updateDesc)
	RootCMD.Flags().BoolVarP(&verbose, "verbose", "v", false, verboseDesc)
	RootCMD.Flags().BoolVar(&version, "version", false, versionDesc)
	RootCMD.Flags().BoolVarP(&tree, "tree", "t", false, treeDesc)
	RootCMD.Flags().BoolVarP(&results, "results", "r", false, resultsDesc)
	RootCMD.Flags().BoolVarP(&lock, "lock", "l", false, lockDesc)
	RootCMD.Flags().BoolVarP(&scan, "scan", "s", false, scanDesc)
	RootCMD.Flags().BoolVarP(&all, "all", "a", false, allDesc)
}

// RootCMD describes the root command.
var RootCMD = &cobra.Command{
	Use:   "govend",
	Short: "Govend vendors external packages.",
	Long:  rootDesc,
	Run: func(cmd *cobra.Command, args []string) {

		if version {
			fmt.Println("0.1.5")
			return
		}

		if scan {
			if err := govend.List(args, format, all); err != nil {
				log.Fatal(err)
			}
			return
		}

		if err := govend.Vendor(args, update, verbose, tree, results, commands, lock, format); err != nil {
			log.Fatal(err)
		}
	},
}
