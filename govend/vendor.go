package govend

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/govend/govend/manifest"
	"github.com/govend/govend/packages"
	"github.com/govend/govend/repo"
	"github.com/govend/govend/semver"
)

var lastimport string

// Vendor
func Vendor(pkgs []string, update, verbose, tree, results, commands, lock bool, format string) error {

	go15, _ := semver.New("1.5.0")
	go16, _ := semver.New("1.6.0")

	version, err := semver.New(strings.TrimPrefix(runtime.Version(), "go"))
	if err != nil {
		return err
	}

	if version.LessThan(go15) {
		return errors.New("govend requires go versions 1.5+")
	}

	if version.GreaterThanEqual(go15) && version.LessThan(go16) {
		if os.Getenv("GO15VENDOREXPERIMENT") != "1" {
			return errors.New("govend requires the env var 'GO15VENDOREXPERIMENT' to be set")
		}
	}

	// attempt to load the manifest file
	m, err := manifest.Load(format)
	if err != nil {
		return err
	}

	// it is important to save the manifest length before syncing, so that
	// if a manifest file existed previously it continues to update even if
	// no current vendors are valid.
	initManifestLen := m.Len()

	// sync ensures that if a vendor is specified in the manifest, that
	// repository root is also currently present in the vendor directory, this
	// allows us to trust the manifest file
	m.Sync()

	// if no packages were provided as arguments, assume the current directory is
	// a go project and scan it for external pacakges.
	if len(pkgs) == 0 {
		pkgs, err = packages.ScanProject(".")
		if err != nil {
			return err
		}
	}

	// download that dependency and any external deps it has
	bpkgs := &badpkgs{}
	numOfpkgs := 0
	for _, pkg := range pkgs {
		lastimport = pkg
		n, err := deptree(pkg, m, bpkgs, 0, verbose, tree)
		if err != nil {
			return err
		}
		numOfpkgs += n
	}

	if verbose && results {
		fmt.Printf("\npackages scanned: %d\n", numOfpkgs)
		fmt.Printf("packages skipped: %d\n", bpkgs.len())
		fmt.Printf("repos downloaded: %d\n", m.Len())
	}

	if lock || initManifestLen > 0 {
		if err := m.Write(); err != nil {
			return err
		}
	}

	return nil
}

// deptree downloads a dependency and the entire tree of dependencies/packages
// that dependency requries as well.
//
// deptree takes a manifest as well as map of badimports to avoid as much
// rework as possible.
//
// as well as an error, deptree returns the number of external package nodes
// scanned in the dependecy tree excluding the root node/pkg.
func deptree(pkg string, m *manifest.Manifest, bpkgs *badpkgs, level int, verbose bool, tree bool) (int, error) {

	// use the network to gather some metadata on this repo
	r, err := repo.Ping(pkg)
	if err != nil {
		if strings.Contains(err.Error(), "unrecognized import path") {
			bpkgs.append(pkg)
			if verbose {
				if tree {
					writeBlanks(level)
				}
				fmt.Printf("%s (bad ping)\n", pkg)
			}
		}
		return 0, err
	}

	// check if the repo is missing from the manifest file
	if !m.Contains(r.ImportPath) {

		if verbose {
			if tree {
				writeBlanks(level)
			}
			fmt.Printf("%s\n", r.ImportPath)
		}

		// download the repo
		rev, err := repo.Download(r, "vendor", "latest")
		if err != nil {
			return 0, err
		}

		// append the repo to the manifest file
		m.Append(r.ImportPath, rev)
	}

	pkgdeps, err := packages.Scan(filepath.Join("vendor", pkg))
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	// exclude standard packages
	pkgdeps = packages.FilterStdPkgs(pkgdeps)
	num := len(pkgdeps)
	level++

	for _, pkgdep := range pkgdeps {
		if pkgdep != pkg && pkgdep != lastimport {

			lastimport = pkg
			n, err := deptree(pkgdep, m, bpkgs, level, verbose, tree)
			if err != nil {
				return num + n, err
			}

			// add num nodes scanned to the total tree n
			num += n
		}
	}

	return num, nil
}

type badpkgs struct {
	pkgs []string
}

func (b *badpkgs) append(pkg string) {
	b.pkgs = append(b.pkgs, pkg)
}

func (b *badpkgs) contains(pkg string) bool {
	for _, p := range b.pkgs {
		if p == pkg {
			return true
		}
	}
	return false
}

func (b *badpkgs) len() int {
	return len(b.pkgs)
}

// writeBlanks writes a number of blank spaces.
func writeBlanks(num int) {
	for num > 0 {
		fmt.Printf(" ")
		num--
	}
}
