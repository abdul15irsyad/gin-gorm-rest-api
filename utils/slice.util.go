package utils

// get random element from array/slice
//
//	random := *RandomSlice([]string{"element1", "element2","element3"})
func RandomSlice[T any](array []T) T {
	var element T
	if len(array) == 0 {
		return element
	}
	return array[RandomInt(0, len(array)-1)]
}

// shuffle an array/slice
//
//	utils.ShuffleSlice([]string{"element1", "element2","element3"})
func ShuffleSlice[T any](array []T) []T {
	for i := range array {
		j := RandomInt(0, len(array)-1)
		array[i], array[j] = array[j], array[i]
	}
	return array
}
