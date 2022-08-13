package slice

func Unique[T comparable](list []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range list {
		if _, exists := seen[item]; !exists {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
