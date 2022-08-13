package slice

func Unique[T comparable](list []T) []T {
	seen := make(map[T]bool)
	uniqueList := []T{}

	for _, item := range list {
		if _, exists := seen[item]; !exists {
			seen[item] = true
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}
