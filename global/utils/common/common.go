package common

// InArray checks if a value exists in a slice
func InArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
