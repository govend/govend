// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

import (
	"fmt"

	"github.com/govend/govend/deps/repos"
	"github.com/govend/govend/manifest"
)

// Hold takes a manifest and downloads all the repos with a hold value of true.
func Hold(m *manifest.Manifest, verbose bool) int {
	var i int
	for _, vendor := range m.Vendors {
		if vendor.Hold {
			if verbose && i == 0 {
				if verbose {
					fmt.Print("\n")
				}
			}
			if verbose {
				fmt.Println(vendor.Path)
			}
			repo, err := repos.Ping(vendor.Path)
			if err != nil {
				fmt.Printf("%s (bad ping): %s\n", vendor.Path, err)
				continue
			}
			rev, err := repos.Download(repo, "vendor", vendor.Rev)
			if err != nil {
				fmt.Printf("%s (download error): %s\n", repo.ImportPath, err)
				continue
			}
			m.Append(vendor.Path, rev, vendor.Hold)
			i++
		}
	}
	return i
}
