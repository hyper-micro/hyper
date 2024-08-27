package slice

func ContainsString(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func Contains[T comparable](v T, sl []T) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
