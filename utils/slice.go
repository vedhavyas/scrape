package utils

//FilterSliceAndAdd filters newList from baseList and return the filteredSlice
func FilterSlice(baseSlice []string, s []string) []string {
	var filteredSlice []string
	for _, data := range s {
		if !IsExists(baseSlice, data) {
			filteredSlice = append(filteredSlice, data)
		}
	}

	return filteredSlice
}

//IsExists will check if the given string exist in given slice
func IsExists(baseSlice []string, data string) bool {
	for _, s := range baseSlice {
		if s == data {
			return true
		}
	}

	return false
}

//RemoveDuplicates will remove duplicate entries from slice
func RemoveDuplicates(s []string) []string {
	m := make(map[string]bool)
	for _, data := range s {
		m[data] = true
	}

	var s2 []string
	for key := range m {
		s2 = append(s2, key)
	}

	return s2
}
