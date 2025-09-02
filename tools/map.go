package tools

func Map2SliceWithHeader(m map[string]any, header []string) []any {
	var slice = make([]any, len(header), len(header))
	for i, k := range header {
		slice[i] = m[k]
	}
	return slice
}
