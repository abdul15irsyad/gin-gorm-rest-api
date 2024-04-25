package utils

func RandomSlice[T any](array []T) *T {
	if len(array) == 0 {
		return nil
	}
	return &array[RandomInt(0, len(array)-1)]
}

// example
// random := *utils.RandomSlice([]string{"title1", "title2","title3"})
