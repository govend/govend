// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kr/fs"
)

// Prune takes a deptree and verbose option and returns the number of
// directories and files removed.
func Prune(deptree []string, verbose bool) (dirs, files int, pruned []string) {
	pruned = make([]string, 0, len(deptree))

	if verbose {
		fmt.Print("\nprune vendored packages... ")
	}

	w := fs.Walk("vendor")
	for w.Step() {

		finfo := w.Stat()
		if finfo.IsDir() {

			if inDepTree(deptree, w.Path()) {
				continue
			}

			os.RemoveAll(w.Path())
			w.SkipDir()
			pruned = append(pruned, w.Path())
			dirs++
			continue
		}

		firstchar := []rune(finfo.Name())[0]
		if firstchar == '_' || firstchar == '.' || strings.HasSuffix(finfo.Name(), "_test.go") {
			os.Remove(w.Path())
			files++
		}
	}

	if verbose {
		fmt.Println("finished!")
	}

	return dirs, files, pruned
}

func inDepTree(deptree []string, path string) bool {
	for _, dep := range deptree {
		if subpath(filepath.Join("vendor", dep), path) {
			return true
		}
	}
	return false
}

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
