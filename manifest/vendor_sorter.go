package manifest

type vendorSorter []Vendor

func (v vendorSorter) Len() int {
	return len(v)
}

func (v vendorSorter) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v vendorSorter) Less(i, j int) bool {
	return v[i].Path < v[j].Path
}
