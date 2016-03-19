// Copyright 2016 govend. All rights reserved.
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
		if finfo.IsDir() {

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
			continue
		}

		firstchar := []rune(finfo.Name())[0]
		if firstchar == '_' || firstchar == '.' || strings.HasSuffix(finfo.Name(), "_test.go") {
			os.Remove(w.Path())
		}
	}
}
