// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

// VendOptions represents available vend options.
type VendOptions int

const (
	// UpdateOption updates vendored repositories.
	UpdateOption VendOptions = iota

	// LockOption locks the revision version of vendored repositories.
	LockOption

	// HoldOption holds onto a vendored repository, even if none of its import
	// paths are used in the project source code.
	HoldOption

	// PruneOption removes vendored packages that are not needed.
	PruneOption

	// IgnoreOption ignores the source import paths.
	IgnoreOption

	// VerboseOption prints out packages as they are vendored.
	VerboseOption

	// TreeOption prints the names of packages as an indented tree.
	TreeOption

	// ResultsOption prints a summary of the number of packages scanned, packages
	// skipped, and repositories downloaded.
	ResultsOption
)

// ParseOptions converts cli flag inputs to VendOptions.
func ParseOptions(update, lock, hold, prune, ignore, verbose, tree, results bool) []VendOptions {
	options := []VendOptions{}
	if update {
		options = append(options, UpdateOption)
	}
	if lock {
		options = append(options, LockOption)
	}
	if hold {
		options = append(options, HoldOption)
	}
	if prune {
		options = append(options, PruneOption)
	}
	if ignore {
		options = append(options, IgnoreOption)
	}
	if verbose {
		options = append(options, VerboseOption)
	}
	if tree {
		options = append(options, TreeOption)
	}
	if results {
		options = append(options, ResultsOption)
	}
	return options
}

func parseVendOptions(options []VendOptions) (update, lock, hold, prune, ignore, verbose, tree, results bool) {
	for _, option := range options {
		switch option {
		case UpdateOption:
			update = true
		case LockOption:
			lock = true
		case HoldOption:
			hold = true
		case PruneOption:
			prune = true
		case IgnoreOption:
			ignore = true
		case VerboseOption:
			verbose = true
		case TreeOption:
			tree = true
		case ResultsOption:
			results = true
		}
	}
	return update, lock, hold, prune, ignore, verbose, tree, results
}
