package repo

import "fmt"

// Ping checks the host server to determine the package Repo.
func Ping(pkg string) (*Repo, error) {

	// determine import path and repository type by asking the server host
	repo, err := ImportPath(pkg, false)
	if err != nil {
		e := err.Error()
		msg := "no go-import meta tags"
		if e[len(e)-len(msg):] == msg {
			return nil, fmt.Errorf("network ping (potential proxy issue): %s ", e)
		}
		return nil, fmt.Errorf("network ping: %s", err)
	}

	return repo, nil
}

/*
func pingrepos(uvpkgs []string, manifest []vendor, localpath, vendorProjectPath string, verbosity bool) (map[string]*RepoRoot, error) {

	// create a repo map of package paths to RepoRoots
	rmap := make(map[string]*RepoRoot)

	// verbosity
	if verbosity {
		fmt.Print("identifying repositories... ")
	}

	// remove package imports that might already be included
	//
	// example: "gopkg.in/mgo.v2/bson" -> "gopkg.in/mgo.v2"
	for _, pkg := range uvpkgs {

		// iterate through the manifest vendors.yml file looking for matches
		// and check if we already have vendored code and a revision number
		for _, v := range manifest {
			if pkg == v.Path {
				if _, err := os.Stat(filepath.Join(localpath, vendorProjectPath, pkg)); err == nil {
					goto UnvendoredPackageMatch
				}
			}
		}

		// check if the package is missing from RepoRoot map
		if _, ok := rmap[pkg]; !ok {

			// determine import path and repository type by asking the server
			// hosting the repository and package
			r, err := RepoRootForImportPath(pkg, false)
			if err != nil {
				e := err.Error()
				msg := "no go-import meta tags"
				if e[len(e)-len(msg):] == msg {
					return nil, errors.New("Are you behind a proxy?" + e + "\n")
				}
				return nil, err
			}

			// if the project package root isn't in the RepoRoot map, add it
			if _, ok := rmap[r.Root]; !ok {
				rmap[r.Root] = r
			}
		}

	UnvendoredPackageMatch:
	}

	// verbosity
	if verbosity {
		fmt.Println("			complete")
	}

	return rmap, nil
}
*/
