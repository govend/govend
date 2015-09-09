package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

var tempVendorDir = "_tmp_vendor"

// Download writes the source contents of a Repo to disk. The revision version
// of the repository is returned and the directory is created.
func Download(r *Repo, dir string) (string, error) {

	if _, err := os.Stat(filepath.Join(dir, r.ImportPath)); err == nil {
		return "", fmt.Errorf("repository '%s' has already been vendored", r.ImportPath)
	}

	if err := r.VCS.Create(filepath.Join(tempVendorDir, r.ImportPath), r.URL); err != nil {
		return "", err
	}

	rev, err := r.VCS.Identify(filepath.Join(tempVendorDir, r.ImportPath))
	if err != nil {
		return "", err
	}

	// mkdir
	if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, r.ImportPath)), 0777); err != nil {
		return "", err
	}

	// copy
	if err := CopyDir(filepath.Join(tempVendorDir, r.ImportPath), filepath.Join(dir, r.ImportPath)); err != nil {
		return "", err
	}

	// remove the temp
	os.RemoveAll(tempVendorDir)

	return rev, nil
}
