package collections

func Unique[Slice ~[]Type, Type comparable](slice Slice) Slice {
	buffer := make(map[Type]struct{}, len(slice))
	for _, elem := range slice {
		buffer[elem] = struct{}{}
	}
	unique := make(Slice, len(buffer))
	i := 0
	for elem := range buffer {
		unique[i] = elem
		i++
	}
	return unique
}

func HasDuplicates[Slice ~[]Type, Type comparable](slice Slice) bool {
	return len(slice) != len(Unique(slice))
}
