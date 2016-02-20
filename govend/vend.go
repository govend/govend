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

// Vend
func Vend(pkgs []string, update, verbose, results, commands, lock bool, format string) error {

	go15, _ := semver.New("1.5.0")
	go16, _ := semver.New("1.6.0")
	go17, _ := semver.New("1.7.0")

	version, err := semver.New(strings.TrimPrefix(runtime.Version(), "go"))
	if err != nil {
		return err
	}

	if version.LessThan(go15) {
		return errors.New("vendoring requires Go versions 1.5+")
	}

	if version.GreaterThanEqual(go15) && version.LessThan(go16) {
		if os.Getenv("GO15VENDOREXPERIMENT") != "1" {
			return errors.New("Go 1.5.x requires 'GO15VENDOREXPERIMENT=1'")
		}
	}

	if version.GreaterThanEqual(go16) && version.LessThan(go17) {
		if os.Getenv("GO15VENDOREXPERIMENT") == "0" {
			return errors.New("Go 1.6.x cannot vendor with 'GO15VENDOREXPERIMENT=0'")
		}
	}

	// attempt to load the manifest file
	m, err := manifest.Load(format)
	if err != nil {
		return err
	}

	// it is important to save the manifest length before syncing, so that
	// we can tell the difference and update the manifest file
	manifestLen := m.Len()

	// sync ensures that if a vendor is specified in the manifest, that the
	// repository structure is also currently present in the vendor directory,
	// this allows us to trust the manifest file
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
	pkglist := map[string]bool{}
	for i := len(pkgs) - 1; i >= 0; i-- {
		if _, ok := pkglist[pkgs[i]]; ok {
			pkgs = pkgs[:i]
			continue
		}
		deps, err := deptree(pkgs[i], m, verbose)
		if err != nil {
			pkglist[pkgs[i]] = false
			pkgs = pkgs[:i]
			continue
		}
		pkglist[pkgs[i]] = true
		pkgs = append(pkgs[:i], pkgs[i+1:]...)
		if len(deps) > 0 {
			pkgs = append(pkgs, deps...)
		}
		i += len(deps)
	}

	if verbose && results {
		fmt.Printf("\npackages scanned: %d\n", len(pkglist))
		fmt.Println("packages skipped:")
		for pkg, ok := range pkglist {
			if !ok {
				fmt.Printf("	%q\n", pkg)
			}
		}
		fmt.Printf("repos downloaded: %d\n", m.Len())
	}

	if lock || manifestLen > 0 {
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
func deptree(pkg string, m *manifest.Manifest, verbose bool) ([]string, error) {

	// use the network to gather some metadata on this repo
	r, err := repo.Ping(pkg)
	if err != nil {
		if strings.Contains(err.Error(), "unrecognized import path") {
			if verbose {
				fmt.Printf("%s (bad ping): %s\n", pkg, err)
			}
		}
		return nil, err
	}

	// check if the repo is missing from the manifest file
	if !m.Contains(r.ImportPath) {
		if verbose {
			fmt.Printf("%s\n", r.ImportPath)
		}
		rev, err := repo.Download(r, "vendor", "latest")
		if err != nil {
			return nil, err
		}

		// append the repo to the manifest file
		m.Append(r.ImportPath, rev)
	}

	pkgdeps, err := packages.Scan(filepath.Join("vendor", pkg))
	if err != nil {
		return nil, err
	}

	// exclude standard packages
	return packages.FilterStdPkgs(pkgdeps), nil
}

// writeBlanks writes a number of blank spaces.
func writeBlanks(num int) {
	for num > 0 {
		fmt.Printf(" ")
		num--
	}
}
