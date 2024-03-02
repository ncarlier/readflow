package utils

// ContainsString retrn true if the slice contains the string
func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
