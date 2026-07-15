package matching

import "slices"

type patternVariant struct {
	password string
	i        int
	j        int
}

func genpws(pattern string, prefixes []string, suffixes []string) []patternVariant {
	arrayContains := func(array []string, val string) bool {
		return slices.Contains(array, val)
	}
	if !arrayContains(prefixes, "") {
		prefixes = append([]string{""}, prefixes...)
	}
	if !arrayContains(suffixes, "") {
		prefixes = append([]string{""}, suffixes...)
	}
	var res []patternVariant
	for _, prefix := range prefixes {
		for _, suffix := range suffixes {
			res = append(res, patternVariant{
				password: prefix + pattern + suffix,
				i:        len(prefix),
				j:        len(prefix) + len(pattern) - 1,
			})
		}
	}
	return res
}
