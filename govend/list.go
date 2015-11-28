package govend

// ScanCMD executes the scan command
import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gophersaurus/govend/packages"
	"github.com/gophersaurus/govend/strutil"
)

// List
func List(args []string, file, format string, all bool, vendors bool) error {

	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// check if any result format has been specified
	if len(format) == 0 {

		// if no file to write to has been specifed default to text,
		// otherwise attempt to determine the file type by file extension
		if len(file) == 0 {
			format = "txt"
		} else {
			ext := path.Ext(file)
			format = ext[1:]
		}
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

	// if a file to write to was specified, write to it
	if len(file) > 0 {
		if err := ioutil.WriteFile(file, b, 0644); err != nil {
			return err
		}
		return nil
	}

	// no file specified, print the results to screen
	fmt.Printf("%s", b)
	return nil
}
