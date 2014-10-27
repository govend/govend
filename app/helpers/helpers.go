package helpers

import (
	"fmt"
	"log"
)

// Check takes an error and fatally logs the results.
func Check(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}
