package tasks

import (
	"fmt"

	"github.com/JackSpirou/govend/app/deps"
	"github.com/JackSpirou/govend/app/helpers"
	"github.com/JackSpirou/govend/app/utils/copyrecur"
)

func copyDeps(gopath string, list bool) {

	if list {
		fmt.Println("")
	}

	// Copy all dependency code into the vendor directory.
	for _, dep := range deps.List {

		// Source and destination paths.
		src := gopath + "/src/" + dep
		dest := deps.Dir + "/" + dep

		if list {
			// Tell the user we are pulling this dependency into their project.
			fmt.Println(" ↓ " + dep)
		}

		// Recursively copy the dependency into the vendors directory.
		err := copyrecur.CopyDir(src, dest)
		helpers.Check(err)
	}

}
