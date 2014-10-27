package tasks

import (
	"errors"
	"fmt"
	"os"

	"github.com/JackSpirou/govend/app/deps"
	"github.com/JackSpirou/govend/app/helpers"
)

func Govend(list, save, update bool) {

	// Read the $GOPATH env.
	gopath := os.Getenv("GOPATH")
	if len(gopath) < 1 {
		helpers.Check(errors.New("env var $GOPATH not found"))
	}

	// Check govend has all its own dependencies installed.
	govendDeps()

	// Pull all dependencies into the $GOPATH.
	goGetDeps(update)

	// Remove the vendor directory if it exists.
	err := os.RemoveAll(deps.Dir)
	helpers.Check(err)

	// Copy the dependencies in the gopath.
	copyDeps(gopath, list)

	// Delete the dependencies in the gopath.
	deleteDeps(gopath)

	// Tell the user we are now attempting to fix imports.
	fmt.Println("")
	fmt.Println("Fixing imports...")

	// Run gormimports
	goRmImports()

	// Run goimports
	goImports()

	// If the user wants to save the dendencies in the gopath, go get them again.
	if save {
		goGetDeps(true)
	}

	// Tell the user we are done vendoring!
	fmt.Println("")
	fmt.Println("Vending complete!")
}
