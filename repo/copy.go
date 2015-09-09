package repo

// Origial code by Jaybill McCarthy: http://jayblog.jaybill.com/post/id/26

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// CopyFile copies a file source to the provided destination.
func CopyFile(source string, dest string) error {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}
	return err
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
//
// CopyDir ignores '.gitignore' files and all other directories that contain a
// '.' or '_' preceding them.
func CopyDir(source string, dest string) error {
	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("Source is not a directory")
	}
	// ensure dest dir does not already exist
	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return errors.New("Destination already exists")
	}
	// create dest dir
	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(source)
	for _, entry := range entries {

		sfp := filepath.Join(source, entry.Name())
		dfp := filepath.Join(dest, entry.Name())
		if entry.IsDir() {
			// Don't copy directories that have a '.' or '_' preceding them
			n := entry.Name()
			if n[0] != '.' && n[0] != '_' {
				err = CopyDir(sfp, dfp)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			if entry.Name() != ".gitignore" {
				// perform copy
				err = CopyFile(sfp, dfp)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	return err
}
