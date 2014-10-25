// Copy a whole directory tree. Current implementation is a packaged version of Jaybill McCarthy's code which can be found at http://jayblog.jaybill.com/post/id/26
package copyrecur

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
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
	return
}

// Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return &CustomError{"Source is not a directory"}
	}
	// ensure dest dir does not already exist
	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return &CustomError{"Destination already exists"}
	}
	// create dest dir
	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(source)
	for _, entry := range entries {
		if entry.Name() != ".git" || entry.Name() != ".gitignore" {
			sfp := source + "/" + entry.Name()
			dfp := dest + "/" + entry.Name()
			if entry.IsDir() {
				err = CopyDir(sfp, dfp)
				if err != nil {
					log.Println(err)
				}
			} else {
				// perform copy
				err = CopyFile(sfp, dfp)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			fmt.Println(entry.Name())
		}

	}
	return
}

// A struct for returning custom error messages
type CustomError struct {
	What string
}

// Returns the error message defined in What as a string
func (e *CustomError) Error() string {
	return e.What
}
