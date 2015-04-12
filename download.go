package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func download(repomap map[string]*RepoRoot, manifest []vendor, vendorTempPath, vendorPath string, verbosity bool) error {

	// verbosity
	if verbosity {
		fmt.Println("downloading packages...")
	}

	// iterate through the repomap
	for key, r := range repomap {

		// create a directory for the pkg
		if err := os.MkdirAll(filepath.Dir(filepath.Join(vendorTempPath, r.Root)), 0777); err != nil {
			return err
		}

		// iterate over the vendors in the vendor file
		for i, v := range manifest {

			// check if we have a match
			if r.Root == v.Path {

				if _, err := os.Stat(filepath.Join(vendorPath, v.Path)); err == nil {

					// verbosity
					if verbosity {
						fmt.Println(" - " + r.Root + " (vendored)")
					}

					delete(repomap, key)

					goto UnvendoredManifestMatch
				}

				// check if a revision exists
				if len(v.Rev) > 0 {

					// verbosity
					if verbosity {
						fmt.Println(" ↓ " + r.Repo + " (" + v.Rev + ")")
					}

					// create the repository at that specific revision
					if err := r.VCS.CreateAtRev(filepath.Join(vendorTempPath, r.Root), r.Repo, v.Rev); err != nil {
						return err
					}
				} else {

					// Create the repository from the default repository revision.
					if err := r.VCS.Create(filepath.Join(vendorTempPath, r.Root), r.Repo); err != nil {
						return err
					}

					rev, err := r.VCS.identify(filepath.Join(vendorTempPath, r.Root))
					if err != nil {
						return err
					}

					manifest[i] = vendor{Path: v.Path, Rev: rev}

					// verbosity
					if verbosity {
						fmt.Println(" ↓ " + r.Repo + " (" + rev + ")")
					}
				}
				goto UnvendoredManifestMatch
			}
		}

		// verbosity
		if verbosity {
			fmt.Println(" ↓ " + r.Repo + " (latest)")
		}

		// Create the repository from the default repository revision.
		if err := r.VCS.Create(filepath.Join(vendorTempPath, r.Root), r.Repo); err != nil {
			return err
		}

		if rev, err := r.VCS.identify(filepath.Join(vendorTempPath, r.Root)); err == nil {
			manifest = append(manifest, vendor{Path: r.Root, Rev: rev})
		} else {
			return err
		}

	UnvendoredManifestMatch:
	}

	// verbosity
	if verbosity {
		fmt.Print("vendoring packages...")
	}

	// iterate through the repomap
	for _, r := range repomap {

		// remove all
		if err := os.RemoveAll(filepath.Join(vendorPath, r.Root)); err != nil {
			return err
		}

		// mkdir
		if err := os.MkdirAll(filepath.Dir(filepath.Join(vendorPath, r.Root)), 0777); err != nil {
			return err
		}

		// copy
		if err := CopyDir(filepath.Join(vendorTempPath, r.Root), filepath.Join(vendorPath, r.Root)); err != nil {
			return err
		}
	}

	// remove the temp
	os.RemoveAll(vendorTempPath)

	// verbosity
	if verbosity {
		fmt.Println("				complete")
	}

	return nil
}
