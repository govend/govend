// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Clean cleans a repository codebase.
func Clean(dir string) error {

	// get properties of source dir
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	// iterate over each entry (file or directory)
	entries, err := ioutil.ReadDir(dir)
	for _, entry := range entries {

		// get the source file path
		sfp := filepath.Join(dir, entry.Name())

		// remove directories that have a '.' or '_' preceding them
		if entry.IsDir() {
			n := entry.Name()
			if n[0] == '.' || n[0] == '_' {
				if err := os.RemoveAll(sfp); err != nil {
					return err
				}
			}
		}

		// remove .gitignore files
		if entry.Name() == ".gitignore" {
			if err := os.Remove(sfp); err != nil {
				return err
			}
		}
	}

	return err
}
