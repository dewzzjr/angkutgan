package pagination

// Offset get offset from page and row
func Offset(page, row int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * row
}
