package util

// SliceContainsString checks if a given string slice s contains
// the target string.
func SliceContainsString(s []string, target string) bool {
	for _, el := range s {
		if target == el {
			return true
		}
	}
	return false
}

// SliceRemoveString removes the target string from slice s.
// If there are duplicates of target present in s, it removes them all.
func SliceRemoveString(s *[]string, target string) {
	result := make([]string, 0)
	for _, el := range *s {
		if el != target {
			result = append(result, el)
		}
	}
	*s = result
}
