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

var semver = "v0.1.11-beta"

// cli flag values
var (
	version       bool
	verbose       bool
	tree          bool
	update        bool
	results       bool
	lock          bool
	hold          bool
	prune         bool
	alltags       bool
	ignore        bool
	scan          bool
	skipTestFiles bool
	skipFilters   bool
	format        string
	strict        bool
	tags          []string
)

// cli flag usage descriptions
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
	alltagsDesc = `The --all-tags flag scans all files, regardless of build
	tags, for imports.
	`
	ignoreDesc = `The --ignore flag ignores any packages that are NOT found in the
	manifest file.
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
	strictDesc = `The --strict flag returns non-zero status code when a package
	path and/or revision is invalid.
	`
	buildtagsDesc = `The --build-tags flag scans files with these build tags
	for imports. It takes a comma-separated list. These should be the tags that
	will actually be used during compilation.
	`
)

// govend represents the command root
var govend = &cobra.Command{
	Use:   "govend",
	Short: "The govend command vendors external packages.",
	Long:  govendDesc,
	Run: func(cmd *cobra.Command, args []string) {

		switch {
		case version:
			fmt.Println(semver)

		case scan:
			// we should always assume the local directory unless a specific
			// directory path is provided
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// parse flag options relevant to the scan command
			scanOpts := []imports.ScanOption{}
			if skipTestFiles {
				scanOpts = append(scanOpts, imports.SkipTestFiles)
			}
			if skipFilters {
				scanOpts = append(scanOpts, imports.SkipFilters)
			}
			if alltags {
				scanOpts = append(scanOpts, imports.AllBuildTags)
			} else {
				scanOpts = append(scanOpts, imports.BuildTags(tags...))
			}

			pkgs, err := imports.Scan(path, scanOpts...)
			if err != nil {
				log.Fatal(err)
			}

			b, err := imports.Format(pkgs, format)
			if err != nil {
				log.Fatal(err)
			}

			// always print the scan results to screen
			fmt.Printf("%s\n", b)

		default:
			// the default govend command triggers vending since this is the most
			// common use case.

			// first we need to check that the current project environment is
			// suitable for vendoring packages, otherwise the user will not get
			// the results they expect when attempting `go build` or `go install`
			if err := deps.Vendorable(verbose); err != nil {
				log.Fatal(err)
			}

			// parse flag options relevant to the vend command
			vOpts := deps.ParseOptions(update, lock, hold, prune, alltags, ignore, verbose, tree, results, strict)

			// vendor according to the options provided
			if err := deps.Vend(args, format, tags, vOpts...); err != nil {
				log.Fatal(err)
			}
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
	govend.Flags().BoolVar(&alltags, "all-tags", false, alltagsDesc)
	govend.Flags().BoolVarP(&ignore, "ignore", "i", false, ignoreDesc)
	govend.Flags().BoolVarP(&scan, "scan", "s", false, scanDesc)
	govend.Flags().BoolVar(&skipFilters, "skipFilters", false, skipFiltersDesc)
	govend.Flags().BoolVar(&skipTestFiles, "skipTestFiles", false, skipTestFilesDesc)
	govend.Flags().BoolVar(&strict, "strict", false, strictDesc)
	govend.Flags().StringSliceVar(&tags, "build-tags", []string{}, buildtagsDesc)
	govend.Execute()
}
