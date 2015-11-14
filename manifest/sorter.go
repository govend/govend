package manifest

type sorter []Vendor

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sorter) Less(i, j int) bool {
	return s[i].Path < s[j].Path
}
