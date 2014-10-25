package app

import (
	"fmt"

	"github.com/jackspirou/govend/app/tasks"
)

// Run the app.
func Run() {
	// Write a newline, just to be pretty.
	fmt.Println("")

	// Right now the only task govend runs is Govend()
	tasks.Govend()
}
