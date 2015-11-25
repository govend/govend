package repo

import "github.com/gophersaurus/govend/vcs"

// Repo RepoRoot
type Repo struct {
	VCS *VCS

	// URL is the repository URL, including scheme
	URL string

	// ImportPath is the import path corresponding to the root of the repository
	ImportPath string
}

// New returns a new Repo.
func New(v *VCS, url, importpath string) *Repo {
	return &Repo{
		VCS:        v,
		URL:        url,
		ImportPath: importpath,
	}
}

// ImportPath returns a new Repo.
func ImportPath(importpath string, verbose bool) (*Repo, error) {
	rr, err := vcs.RepoRootForImportPath(importpath, verbose)
	if err != nil {
		return nil, err
	}
	vcs, err := NewVCS(rr.VCS)
	if err != nil {
		return nil, err
	}
	return New(vcs, rr.Repo, rr.Root), nil
}

// ImportDynamic returns a new Repo.
func ImportDynamic(importpath string, verbose bool) (*Repo, error) {
	rr, err := vcs.RepoRootForImportDynamic(importpath, verbose)
	if err != nil {
		return nil, err
	}
	vcs, err := NewVCS(rr.VCS)
	if err != nil {
		return nil, err
	}
	return New(vcs, rr.Repo, rr.Root), nil
}
