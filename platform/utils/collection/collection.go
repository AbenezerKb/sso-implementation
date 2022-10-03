package collection

func Contains[T comparable](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func ContainsWithMatcher[T comparable, B any](value T, array []B, matcher func(value T, arrayValue B) bool) bool {
	for i := 0; i < len(array); i++ {
		if matcher(value, array[i]) {
			return true
		}
	}

	return false
}
