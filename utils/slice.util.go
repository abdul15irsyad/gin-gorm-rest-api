package utils

// get random element from array/slice
//
//	random := *RandomArray([]string{"element1", "element2","element3"})
func RandomArray[T any](array []T) *T {
	if len(array) == 0 {
		return nil
	}
	return &array[RandomInt(0, len(array)-1)]
}

// shuffle an array/slice
//
//	utils.ShuffleArray([]string{"element1", "element2","element3"})
func ShuffleArray[T any](array []T) []T {
	for i := range array {
		j := RandomInt(0, len(array)-1)
		array[i], array[j] = array[j], array[i]
	}
	return array
}
