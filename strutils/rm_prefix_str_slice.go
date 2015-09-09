package strutils

// RemovePrefixInStringSlice removes any items from a given slice of strings
// that start with the provided prefix.
func RemovePrefixInStringSlice(prefix string, items []string) []string {

	// determine the length of prefix, we only want to do this once
	l := len(prefix)

	// create an empty slice of results to fill.
	var results []string

	// iterate through the slice of items
	for _, item := range items {

		// check the item length is geater than or equal to the prefix length
		// this ensures no out of bounds memory errors will occur
		if len(item) >= l {

			// check if the prefix matchs the current item
			if prefix == item[:l] {
				continue
			}
		}

		// append the item to the slice
		results = append(results, item)
	}

	return results
}
