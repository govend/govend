package tasks

import (
	"os"

	"github.com/JackSpirou/govend/app/deps"
	"github.com/JackSpirou/govend/app/helpers"
)

func deleteDeps(gopath string) {

	// Delete all files go get downloaded.
	for _, dep := range deps.List {

		src := gopath + "/src/" + dep

		err := os.RemoveAll(src)
		helpers.Check(err)
	}
}
