package collections

import (
	"iter"
)

func Iter[S ~[]T, T any](slice S) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, elem := range slice {
			if !yield(elem) {
				return
			}
		}
	}
}

func Filter[T any](seq iter.Seq[T], by func(T) bool) (it iter.Seq[T], cnt int) {
	it = filter(seq, by)
	for range it {
		cnt++
	}
	return it, cnt
}

func filter[T any](seq iter.Seq[T], by func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range seq {
			if by(i) {
				if !yield(i) {
					return
				}
			}
		}
	}
}

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
