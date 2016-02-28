// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/govend/govend/deps"
	"github.com/govend/govend/imports"
	"github.com/spf13/cobra"
)

var (
	version       bool
	verbose       bool
	update        bool
	results       bool
	lock          bool
	scan          bool
	skipTestFiles bool
	skipFilters   bool
	format        string
)

const (
	govendDesc = `The govend command scans and downloads dependencies.
	`
	versionDesc = `The --version flag prints the installed version of govend.
	`
	verboseDesc = `The -v flag prints package paths as they are vendored.
	`
	resultsDesc = `The -r flag works with the -v flag to print a summary of the
	number of packages scanned, packages skipped, and repositories downloaded.
	`
	formatDesc = `The -f flag works with the -l flag and -s flag to define a
	format. The default format is YAML, but JSON and TOML is also supported.
	`
	updateDesc = `The -u flag uses the network to update packages and their
	dependencies. By default the network is used to check out missing
	packages and ensure correct revision versions.
	`
	lockDesc = `The -l flag locks down dependency versions by writing to disk a
	manifest vendor file containing repository revision hashes.
	`
	scanDesc = `The -s flag scans the current or provided directory for external
	packages.
	`
	skipFiltersDesc = `The --skipFilters flag works with the -s flag to show the
	raw unfiltered list of scanned packages.
	`
	skipTestFilesDesc = `The --skipTestFiles flag works with the -s flag and
	default govend command to skip scanning files that end in "_test.go".
	`
)

// govend describes the root command.
var govend = &cobra.Command{
	Use:   "govend",
	Short: "Govend vendors external packages.",
	Long:  govendDesc,
	Run: func(cmd *cobra.Command, args []string) {

		// print version
		if version {
			fmt.Println("0.1.5")
			return
		}

		// scan for dependencies
		scanOptions := parseScanOptions(skipTestFiles, skipFilters)
		if scan {

			// assume local directory
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// scan the project directory provided
			pkgs, err := imports.Scan(path, scanOptions...)
			if err != nil {
				log.Fatal(err)
			}

			b, err := imports.Format(pkgs, format)
			if err != nil {
				log.Fatal(err)
			}

			// print the results to screen
			fmt.Printf("%s\n", b)
			return
		}

		// vendor dependencies
		if err := deps.Vend(args, update, verbose, results, lock, format); err != nil {
			log.Fatal(err)
		}
	},
}

func parseScanOptions(skipTestFiles, skipFilters bool) []imports.ScanOptions {
	options := []imports.ScanOptions{}
	if skipTestFiles {
		options = append(options, imports.SkipTestFiles)
	}
	if skipFilters {
		options = append(options, imports.SkipFilters)
	}
	return options
}

func main() {
	govend.Flags().BoolVar(&version, "version", false, versionDesc)
	govend.Flags().BoolVarP(&verbose, "verbose", "v", false, verboseDesc)
	govend.Flags().BoolVarP(&results, "results", "r", false, resultsDesc)
	govend.Flags().StringVarP(&format, "format", "f", "YAML", formatDesc)
	govend.Flags().BoolVarP(&update, "update", "u", false, updateDesc)
	govend.Flags().BoolVarP(&lock, "lock", "l", false, lockDesc)
	govend.Flags().BoolVarP(&scan, "scan", "s", false, scanDesc)
	govend.Flags().BoolVar(&skipFilters, "skipFilters", false, skipFiltersDesc)
	govend.Flags().BoolVar(&skipTestFiles, "skipTestFiles", false, skipTestFilesDesc)
	govend.Execute()
}
