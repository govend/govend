// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kr/fs"
)

func subpath(base, path string) bool {
	b := strings.Split(base, "/")
	p := strings.Split(path, "/")
	for i, v := range b {
		if i+1 < len(p) && v == p[i] {
			continue
		}
		return v == p[i] || p[i] == ""
	}
	return false
}

func prunePackages(keepers []string) {
	w := fs.Walk("vendor")
	for w.Step() {
	START:

		finfo := w.Stat()

		// check for any walker errors
		// if w.Err() != nil {
		// log.Fatal(w.Err())
		// }

		if !finfo.IsDir() {
			continue
		}

		for _, keeper := range keepers {
			if subpath(filepath.Join("vendor", keeper), w.Path()) {
				if !w.Step() {
					return
				}
				goto START
			}
		}

		os.RemoveAll(w.Path())
		w.SkipDir()
	}
}
