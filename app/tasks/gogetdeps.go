package tasks

import (
	"fmt"
	"os/exec"

	"github.com/JackSpirou/govend/app/deps"
	"github.com/JackSpirou/govend/app/helpers"
)

func goGetDeps(update bool) {

	fmt.Println("")

	if update {
		fmt.Println("Downloading and checking dependency updates...")
	} else {
		fmt.Println("Downloading dependencies...")
	}

	// Run "$ go get" for all deps.
	for _, dep := range deps.List {

		// Check if updates are needed.
		if update {
			err := exec.Command("go", "get", "-u", "-d", dep).Run()
			helpers.Check(err)
		} else {
			err := exec.Command("go", "get", "-d", dep).Run()
			helpers.Check(err)
		}
	}

}
