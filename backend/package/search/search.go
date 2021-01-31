package search

// ArrayString search string from array string
func ArrayString(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
