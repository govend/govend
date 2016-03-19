// Copyright 2016 govend. All rights reserved.
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
	tree          bool
	update        bool
	results       bool
	lock          bool
	hold          bool
	prune         bool
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
	treeDesc = `The -t flag works with the -v flag to print the names of packages
	as an indented tree.
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
	holdDesc = `The --hold flag holds on to a dependency, even if it is not used
	as an import path in the project codebase. This ability to freeze
	dependencies is useful for vendoring Go tool versions per project.
	`
	pruneDesc = `The --prune flag removes vendored packages that are not needed
	by leveraging the dependency tree after vendoring has completed.
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

		// print version command
		if version {
			fmt.Println("0.1.5")
			return
		}

		// scan command
		if scan {
			// assume local directory
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// scan the project directory provided
			scanOptions := parseScanOptions(skipTestFiles, skipFilters)
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

		// all that is left is the vendor command, but first we need to check that
		// the current project environment is vendorable
		if err := deps.Vendorable(); err != nil {
			log.Fatal(err)
		}

		// vendor dependencies
		vendOptions := parseVendOptions(update, lock, hold, prune, verbose, tree, results)
		if err := deps.Vend(args, format, vendOptions...); err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	govend.Flags().BoolVar(&version, "version", false, versionDesc)
	govend.Flags().BoolVarP(&verbose, "verbose", "v", false, verboseDesc)
	govend.Flags().BoolVarP(&tree, "tree", "t", false, treeDesc)
	govend.Flags().BoolVarP(&results, "results", "r", false, resultsDesc)
	govend.Flags().StringVarP(&format, "format", "f", "YAML", formatDesc)
	govend.Flags().BoolVarP(&update, "update", "u", false, updateDesc)
	govend.Flags().BoolVarP(&lock, "lock", "l", false, lockDesc)
	govend.Flags().BoolVar(&hold, "hold", false, holdDesc)
	govend.Flags().BoolVar(&prune, "prune", false, pruneDesc)
	govend.Flags().BoolVarP(&scan, "scan", "s", false, scanDesc)
	govend.Flags().BoolVar(&skipFilters, "skipFilters", false, skipFiltersDesc)
	govend.Flags().BoolVar(&skipTestFiles, "skipTestFiles", false, skipTestFilesDesc)
	govend.Execute()
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

func parseVendOptions(update, lock, hold, prune, verbose, tree, results bool) []deps.VendOptions {
	options := []deps.VendOptions{}
	if update {
		options = append(options, deps.Update)
	}
	if lock {
		options = append(options, deps.Lock)
	}
	if hold {
		options = append(options, deps.Hold)
	}
	if prune {
		options = append(options, deps.Prune)
	}
	if verbose {
		options = append(options, deps.Verbose)
	}
	if tree {
		options = append(options, deps.Tree)
	}
	if results {
		options = append(options, deps.Results)
	}
	return options
}
