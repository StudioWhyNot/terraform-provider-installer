package system

func WrapString(str string, wrapper string) string {
	return wrapper + str + wrapper
}

// Performs a left-join to merge the rhs map into the lhs map into a new map.
func MergeMaps[T comparable](lhs, rhs map[T]T) map[T]T {
	lhsCLone := CloneMap(lhs)
	if rhs == nil {
		return lhsCLone
	}
	for k, v := range rhs {
		lhsCLone[k] = v
	}
	return lhsCLone
}

func CloneMap[T comparable](val map[T]T) map[T]T {
	if val == nil {
		return map[T]T{}
	}
	newVal := make(map[T]T)
	for k, v := range val {
		newVal[k] = v
	}
	return newVal
}
