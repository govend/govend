package tasks

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/jackspirou/govend/app/deps"
	"github.com/jackspirou/govend/app/helpers"
	"github.com/jackspirou/govend/app/utils/copyrecur"
)

func Govend() {

	fmt.Println("go get rmimports")

	// deps for govend.
	err := exec.Command("go", "get", "-d", "github.com/jackspirou/rmimports").Run()
	helpers.Check(err)

	fmt.Println("go get gormimports")

	// deps for govend.
	err = exec.Command("go", "get", "github.com/jackspirou/gormimports").Run()
	helpers.Check(err)

	fmt.Println("after gormimports")

	// Run "$ go get" for all deps.
	for _, dep := range deps.List {
		// Command to go get {dep/repo}
		fmt.Println(dep)
		err := exec.Command("go", "get", "-d", dep).Run()
		helpers.Check(err)
	}

	// Read the $GOPATH env.
	gopath := os.Getenv("GOPATH")
	if len(gopath) < 1 {
		log.Fatalln("Unable to read the $GOPATH environment variable.")
	}

	// Remove the vendor directory if it exists.
	err = os.RemoveAll(deps.Dir)
	helpers.Check(err)

	// Copy all dependency code into the vendor directory.
	for _, dep := range deps.List {

		// Source and destination paths.
		src := gopath + "/src/" + dep
		dest := deps.Dir + "/" + dep

		// Tell the user we are pulling this dependency into their project.
		fmt.Println(" ↓ " + dep)

		// Recursively copy the dependency into the vendors directory.
		err := copyrecur.CopyDir(src, dest)
		helpers.Check(err)
	}

	// Tell the user we are now attempting to fix imports.
	fmt.Println("")
	fmt.Println("Attempting to fix imports...")

	// Delete all files go get downloaded.
	for _, dep := range deps.List {

		src := gopath + "/src/" + dep

		fmt.Println("Removing: " + src)

		err := os.RemoveAll(src)
		helpers.Check(err)
	}

	fmt.Println("running gormimports")
	// Run gormimports
	err = exec.Command("bash", "-c", "gormimports -w ./*").Run()
	fmt.Println(err)

	fmt.Println("running goimports")
	// Run goimports
	err = exec.Command("bash", "-c", "goimports -w ./*").Run()
	fmt.Println(err)

	// Tell the user we are done vendoring!
	fmt.Println("")
	fmt.Println("Vending complete")

}
