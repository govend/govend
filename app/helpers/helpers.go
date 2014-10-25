package helpers

import "log"

// Check takes an error and fatally logs the results.
func Check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
