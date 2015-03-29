package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/codegangsta/cli"
)

// USAGE
//
// COMMANDS: govend (maybe vend?), scan, prune, update, remove
//
// govend [global options..]
//        -v, -verbose
//

const (
	vendorBase    = "internal"
	vendorDir     = "_vendor"
	vendorFile    = "vendors.yml"
	vendorTempDir = "_tmp_vendor"
)

var (
	vendorPath     = filepath.Join(vendorBase, vendorDir)
	vendorTempPath = filepath.Join(vendorBase, vendorTempDir)
	vendorFilePath = filepath.Join(vendorPath, vendorFile)
)

func main() {

	// Limit go procs to number of CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// start a new cli application
	app := cli.NewApp()

	// define the cli application metadata
	app.Name = "govend"
	app.Usage = "A CLI tool for vendoring golang packages."
	app.Version = "0.0.1"
	app.Author = "Jack Spirou"
	app.Email = "jack.spirou@me.com"

	// define the list of global options.
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "print things as they happen",
		},
	}

	// define the list of commands.
	app.Commands = []cli.Command{
		{
			Name:        "scan",
			Usage:       "Scans a go project for external package dependencies",
			ShortName:   "s",
			Description: "Use this command to find external package dependencies.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "write, w",
					Usage: "write the results to disk",
				},
				cli.StringFlag{
					Name:  "fmt, f",
					Usage: "format the results with values json, yaml, yml, or xml",
				},
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "show all packages, even those in the standard library",
				},
			},
			Action: func(c *cli.Context) {
				if err := scancmd(c.Args().First(), c.String("write"), c.String("fmt"), c.Bool("all")); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:        "imports",
			Usage:       "Rewrites imports prioritizing the projects vendor directory",
			ShortName:   "i",
			Description: "Use this command to for goimports functionality that prioritizes the projects vendor directory imports.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "d",
					Usage: "display diffs instead of rewriting files",
				},
				cli.BoolFlag{
					Name:  "e",
					Usage: "report all errors (not just the first 10 on different lines)",
				},
				cli.BoolFlag{
					Name:  "l",
					Usage: "list files whose formatting differs from goimport's",
				},
				cli.StringFlag{
					Name:  "p",
					Usage: "comma seperated import path prefixes that if matched, take priority",
					Value: "",
				},
				cli.BoolFlag{
					Name:  "w",
					Usage: "write result to (source) file instead of stdout",
				},
			},
			Action: func(c *cli.Context) {
				if err := importcmd(c.Bool("d"), c.Bool("e"), c.Bool("l"), c.Bool("w"), c.String("p"), c.Args()); err != nil {
					panic(err)
				}
			},
		},
	}

	// define the default action.
	app.Action = func(c *cli.Context) {
		if err := vendcmd(c.GlobalBool("verbose")); err != nil {
			panic(err)
		}
	}

	// execute the cli command given.
	app.Run(os.Args)
}
