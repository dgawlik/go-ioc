package goioc

func find[T any](slice []T, pred func(T) bool) int {
	for i, elem := range slice {
		if pred(elem) {
			return i
		}
	}

	return -1
}
