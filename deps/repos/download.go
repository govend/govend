// Copyright 2016 govend. All rights reserved.
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

	downloadDir := filepath.Join(dir, r.ImportPath)

	// if that repo already exists, remove it
	if err := os.RemoveAll(downloadDir); err != nil {
		return "", err
	}

	// create parents of the download dir, as Create() and CreateAtRev()
	// expect it to exist
	if err := os.MkdirAll(filepath.Dir(downloadDir), 0777); err != nil {
		return "", err
	}

	// create the repository at a specific revision
	if vendorRev == "" || vendorRev == "latest" {
		if err := r.VCS.Create(downloadDir, r.URL); err != nil {
			return "", err
		}
	} else {
		if err := r.VCS.CreateAtRev(downloadDir, r.URL, vendorRev); err != nil {
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
