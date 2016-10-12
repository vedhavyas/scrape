package filters

//FilterSliceAndAdd filters newList from baseList and return the filteredSlice
func FilterSlice(baseList []string, newList []string) []string {
	var filteredSlice []string

	for _, newItem := range newList {
		var exists bool
	InnerLoop:
		for _, baseItem := range baseList {
			if baseItem == newItem {
				exists = true
				break InnerLoop
			}
		}

		if !exists {
			filteredSlice = append(filteredSlice, newItem)
		}
	}

	return filteredSlice
}
