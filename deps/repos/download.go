// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package repos

import (
	"os"
	"path/filepath"
)

// Download writes the source contents of a Repo to disk. The revision version
// of the repository is returned and the directory is created.
func Download(r *Repo, dir, vendorRev string) (string, error) {

	// if that repo already exists, remove it
	if err := os.RemoveAll(filepath.Join(dir, r.ImportPath)); err != nil {
		return "", err
	}

	// create the repository at a specific revision
	if vendorRev == "" || vendorRev == "latest" {
		if err := r.VCS.Create(filepath.Join(dir, r.ImportPath), r.URL); err != nil {
			return "", err
		}
	} else {
		if err := r.VCS.CreateAtRev(filepath.Join(dir, r.ImportPath), r.URL, vendorRev); err != nil {
			return "", err
		}
	}

	// identify the particular revision of the codebase
	rev, err := r.VCS.Identify(filepath.Join(dir, r.ImportPath))
	if err != nil {
		return "", err
	}

	if err := Clean(filepath.Join(dir, r.ImportPath)); err != nil {
		return "", err
	}

	return rev, nil
}
