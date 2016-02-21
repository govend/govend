package govend

import (
	"fmt"

	"github.com/govend/govend/packages"
)

// Scan scans and prints external imported packages.
func Scan(args []string, format string, testfiles, all bool) error {

	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// scan the project directory provided
	pkgs, err := packages.Scan(path, false, testfiles, all)
	if err != nil {
		return err
	}

	b, err := packages.Format(pkgs, format)
	if err != nil {
		return err
	}

	// print the results to screen
	fmt.Printf("%s\n", b)
	return nil
}
