package filters

//FilterSliceAndAdd filters newList from baseList and return the filteredSlice
func FilterSlice(baseSlice []string, newSlice []string) []string {
	var filteredSlice []string
	for _, newItem := range newSlice {
		if !IsExists(baseSlice, newItem) {
			filteredSlice = append(filteredSlice, newItem)
		}
	}

	return filteredSlice
}

func IsExists(baseSlice []string, data string) bool {
	for _, s := range baseSlice {
		if s == data {
			return true
		}
	}

	return false
}
