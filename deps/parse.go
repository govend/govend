// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

// ParseVendOptions converts cli flag inputs to VendOptions.
func ParseVendOptions(update, lock, hold, prune, verbose, tree, results bool) []VendOptions {
	options := []VendOptions{}
	if update {
		options = append(options, Update)
	}
	if lock {
		options = append(options, Lock)
	}
	if hold {
		options = append(options, Hold)
	}
	if prune {
		options = append(options, Prune)
	}
	if verbose {
		options = append(options, Verbose)
	}
	if tree {
		options = append(options, Tree)
	}
	if results {
		options = append(options, Results)
	}
	return options
}
