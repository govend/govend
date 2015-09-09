package main

import "github.com/gophersaurus/govend/commands"

// USAGE
//
// COMMANDS: govend (maybe vend?), scan, prune, update, remove
//
// govend [global options..]
//        -v, -verbose
//

/*
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
*/

func main() {
	commands.RootCMD.AddCommand(commands.ListCMD)
	commands.RootCMD.AddCommand(commands.VersionCMD)
	commands.RootCMD.AddCommand(commands.ImportsCMD)
	commands.RootCMD.Execute()
}

/*

func main() {


	// define the cli application metadata
	app.Name = "govend"
	app.Usage = "A CLI tool for vendoring golang packages."
	app.Version = "0.0.1"
	app.Author = "github.com/jackspirou"


	// define the list of commands.
	app.Commands = []cli.Command{
		{
			Name:        "vendor",
			Usage:       "Vendors a go project's external package dependencies",
			ShortName:   "vend",
			Description: "Use this command to vendor external package dependencies.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "recursive, r",
					Usage: "recursively vendor external package dependencies",
				},
			},
			Action: func(c *cli.Context) {
				if err := vendcmd(c.GlobalBool("verbose"), c.Bool("recursive")); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:        "update",
			Usage:       "Update updates external vendored go packages",
			ShortName:   "u",
			Description: "Use this command to update external vendored package dependencies.",
			Action: func(c *cli.Context) {
				if err := updatecmd(c.Args().First(), c.Args().Get(2), c.GlobalBool("verbose"), c.Bool("recursive")); err != nil {
					panic(err)
				}
			},
		},
	}

	// define the default action.
	app.Action = func(c *cli.Context) {
		if err := vendcmd(c.GlobalBool("verbose"), false); err != nil {
			panic(err)
		}
	}

	// execute the cli command given.
	app.Run(os.Args)
}
*/
