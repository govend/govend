package main

import "github.com/gophersaurus/govend/commands"

func main() {
	commands.RootCMD.AddCommand(commands.ListCMD)
	commands.RootCMD.AddCommand(commands.VersionCMD)
	commands.RootCMD.AddCommand(commands.ImportsCMD)
	commands.RootCMD.Execute()
}
