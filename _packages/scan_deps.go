package packages

import "github.com/gophersaurus/govend/strutil"

// ScanDeps
func ScanDeps(dir, vendorDir string, skipVendor, verbose bool) ([]string, map[string]string, error) {

	// scan for external packages
	pkgs, badpkgs, err := Scan(dir, vendorDir, skipVendor, verbose)
	if err != nil {
		return nil, badpkgs, err
	}

	// filter out standard packages
	pkgs = RemoveStandardPackages(pkgs)

	// filter out internal packages
	projectImportPath, err := ImportPath(dir)
	if err != nil {
		return nil, badpkgs, err
	}
	pkgs = strutil.RemovePrefixInStringSlice(projectImportPath, pkgs)

	return pkgs, badpkgs, nil
}
