package collections

import (
	"iter"
	"slices"
)

func Unique[Slice ~[]Type, Type comparable](slice Slice) Slice {

	buffer := make(map[Type]struct{})
	for _, elem := range slice {
		buffer[elem] = struct{}{}
	}

	uniqueIter := func() iter.Seq[Type] {
		return func(yield func(Type) bool) {
			for elem := range buffer {
				if !yield(elem) {
					return
				}
			}
		}
	}
	return slices.Collect(uniqueIter())
}

func Unique2[Slice ~[]Type, Type comparable](slice Slice) Slice {
	buffer := make(map[Type]struct{})
	for _, elem := range slice {
		buffer[elem] = struct{}{}
	}
	unique := make(Slice, 0, len(buffer))
	for elem := range buffer {
		unique = append(unique, elem)
	}
	return unique
}

func HasDuplicates[Slice ~[]Type, Type comparable](slice Slice) bool {
	return len(slice) != len(Unique(slice))
}
