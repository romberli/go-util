package common

func StringInSlice(str string, slice []string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}

	return false
}

func StringInMap(str string, m map[string]string) bool {
	if _, ok := m[str]; ok {
		return true
	}

	return false
}
