// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package imports

import (
	"errors"
	"os"
	"path/filepath"
)

// Path returns the import path of a package directory.
// It does so via $GOPATH.
func Path(dir string) (string, error) {

	// check for relative path
	if dir == "." {
		dir = os.Args[0]
	}

	// determine the absolute file path for the provided directory
	path, err := filepath.Abs(filepath.Dir(dir))
	if err != nil {
		return "", err
	}

	// check the env $GOPATH is valid
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		return "", errors.New("$GOPATH not set")
	}

	// leverage the $GOPATH to strip out everything but the base git URL
	return path[len(gopath+"/src/"):], nil
}
