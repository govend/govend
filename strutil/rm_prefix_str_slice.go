package strutil

// RemovePrefixInStringSlice removes any items from a given slice of strings
// that start with the provided prefix.
func RemovePrefixInStringSlice(prefix string, items []string) []string {

	l := len(prefix)

	var results []string

	// iterate through the slice of items
	for _, item := range items {

		// check the item length is geater than or equal to the prefix length
		// this ensures no out of bounds memory errors will occur
		if len(item) >= l {
			if prefix == item[:l] {
				continue
			}
		}

		results = append(results, item)
	}

	return results
}
