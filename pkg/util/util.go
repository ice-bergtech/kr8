package util

// Filter returns a new slice containing only the elements that satisfy the predicate function.
// From https://gobyexample.com/collection-functions
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
