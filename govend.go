package main

//go:generate go run mkpkgs.go

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

// USAGE
//
// govend [global options..]
//        -v, -verbose
//

func main() {

	// do the max prox runtime CPU thang.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// start a new cli application.
	app := cli.NewApp()

	// define the cli application metadata.
	app.Name = "govend"
	app.Usage = "A CLI tool for vendoring golang packages."
	app.Version = "0.0.1"
	app.Author = "Jack Spirou"
	app.Email = "jack.spirou@me.com"

	// define the list of global options.
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Verbose Mode: Print things as they happen.",
		},
	}

	// define the list of commands.
	app.Commands = []cli.Command{}

	// define the default action.
	app.Action = func(c *cli.Context) {
		if err := vend(c.GlobalBool("verbose")); err != nil {
			panic(err)
		}
	}

	// execute the cli command given.
	app.Run(os.Args)
}
