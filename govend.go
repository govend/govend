package main

import (
	"os"
	"runtime"

	"github.com/JackSpirou/govend/app"
)

// The go program starts here.
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Run the app.
	exitCode := app.Run()

	// Exit with the proper code.
	os.Exit(exitCode)
}
