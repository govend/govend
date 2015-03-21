package tasks

import (
	"os"

	"github.com/jackspirou/govend/app/deps"
	"github.com/jackspirou/govend/app/helpers"
)

func deleteDeps(gopath string) {

	// Delete all files go get downloaded.
	for _, dep := range deps.List {

		src := gopath + "/src/" + dep

		err := os.RemoveAll(src)
		helpers.Check(err)
	}
}
