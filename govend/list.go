package govend

// ScanCMD executes the scan command
import (
	"fmt"

	"github.com/govend/govend/packages"
	"github.com/govend/govend/strutil"
)

// List
func List(args []string, format string, all bool) error {

	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// scan the project directory provided
	pkgs, err := packages.Scan(dir)
	if err != nil {
		return err
	}

	// remove standard packages
	if !all {
		pkgs = packages.FilterStdPkgs(pkgs)
	}

	projectpath, err := packages.ImportPath(dir)
	if err != nil {
		return err
	}

	// filter out packages internal to the project
	pkgs = strutil.RemovePrefixInStringSlice(projectpath, pkgs)

	b, err := packages.Format(pkgs, format)
	if err != nil {
		return err
	}

	// print the results to screen
	fmt.Printf("%s\n", b)
	return nil
}
