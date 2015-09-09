package packages

import "github.com/gophersaurus/govend/strutils"

// ScanDeps
func ScanDeps(dir, vendorDir string, skipVendor bool) ([]string, error) {

	// scan for external packages
	pkgs, err := Scan(dir, vendorDir, skipVendor)
	if err != nil {
		return nil, err
	}

	// filter out standard packages
	pkgs = RemoveStandardPackages(pkgs)

	// filter out internal packages
	projectImportPath, err := ImportPath(dir)
	if err != nil {
		return nil, err
	}
	pkgs = strutils.RemovePrefixInStringSlice(projectImportPath, pkgs)

	return pkgs, nil
}
