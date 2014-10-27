package app

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/JackSpirou/govend/app/helpers"
	"github.com/JackSpirou/govend/app/tasks"
)

var (
	// flag modes
	list    = flag.Bool("l", false, "list the dependency paths for the project")
	save    = flag.Bool("s", false, "save the project dependencies also to the $GOPATH src")
	update  = flag.Bool("u", false, "update the project dependencies")
	version = flag.Bool("v", false, "display govend's current version number")
)

// Run the app.
func Run() int {

	// Declare usage and parse flags.
	flag.Usage = usage
	flag.Parse()

	// Validate flags.
	if flag.NFlag() > 4 {
		helpers.Report(errors.New("too many flags"))
		usage()
	}

	// Check for version flag.
	if *version {

		// Report version and get exit code.
		exitCode := tasks.Version()

		// Check if user wanted other options... cause thats not ok.
		if *list || *save || *update {
			helpers.Report(errors.New("Don't run -v with other flags."))
		}

		return exitCode
	}

	// Right now the only task govend runs is Govend()
	tasks.Govend(*list, *save, *update)
	return 0
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: govend [flags]\n")
	flag.PrintDefaults()
	os.Exit(2)
}
