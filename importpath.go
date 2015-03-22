package main

import (
	"errors"
	"os"
	"path/filepath"
)

// importpath returns an import path by using $GOPATH.
func importpath(dir string) (string, error) {

	// check if a relative path was provided
	if dir == "." {
		dir = os.Args[0]
	}

	// determine the absolute file path for the provided directory
	path, err := filepath.Abs(filepath.Dir(dir))

	// Check for errors.
	if err != nil {
		return "", err
	}

	// Identify the local $GOPATH.
	gopath := os.Getenv("GOPATH")

	// Check the $GOPATH value is valid.
	if len(gopath) == 0 {
		return "", errors.New("$GOPATH not set")
	}

	// Use the $GOPATH to strip everything out, but the base git URL.
	return path[len(gopath+"/src/"):], nil
}
