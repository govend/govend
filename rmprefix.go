package main

// rmprefix removes any items from a string slice that start with that prefix
func rmprefix(prefix string, items []string) []string {

	// determine the length of prefix, we only want to do this once
	l := len(prefix)

	// iterate through the slice of items
	for i, item := range items {

		// check the prefix length is less than or equal to the next item
		// this ensures no out of bounds memory errors will occur
		if l <= len(item) {

			// check if the prefix matches the current item
			if item[:l] == prefix {

				// remove that item from the slice
				items = append(items[:i], items[i+1:]...)
			}
		}
	}

	return items
}
