package main

// rmprefix removes any items from a string slice that start with that prefix
func rmprefix(prefix string, items []string) []string {

	// determine the length of prefix, we only want to do this once
	l := len(prefix)

	// create an empty slice of results to fill.
	var results []string

	// iterate through the slice of items
	for _, item := range items {

		// check the prefix length is less than or equal to the next item
		// this ensures no out of bounds memory errors will occur
		if l <= len(item) {

			// check if the prefix matchs the current item
			if prefix == item[:l] {
				continue
			}
		}

		// remove that item from the slice
		results = append(results, item)
	}

	return results
}
