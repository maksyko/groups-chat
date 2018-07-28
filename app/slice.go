package app

func SliceContains(slice []string, elem string) bool {
	for _, t := range slice {
		if t == elem {
			return true
		}
	}
	return false
}

func RemoveFromSlice(slice []string, elem string) []string {
	idx := -1
	for i, el := range slice {
		if el == elem {
			idx = i
			break
		}
	}
	if idx > -1 {
		slice = append(slice[:idx], slice[idx+1:]...)
	}
	return slice
}
