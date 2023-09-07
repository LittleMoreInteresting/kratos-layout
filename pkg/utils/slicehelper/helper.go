package slicehelper

import (
	"sort"
)

type Number interface {
	int | int8 | int32 | int64 | float32 | float64
}

func InSliceInt64(tag int64, array []int64) bool {

	for _, n := range array {
		if tag == n {
			return true
		}
	}
	return false
}

func InSlice[T Number](tag T, slice []T) bool {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	search := sort.Search(len(slice), func(i int) bool {
		return slice[i] >= tag
	})
	if search < len(slice) && slice[search] == tag {
		return true
	}
	return false
}

func SimpleMergeNumber[T Number](s ...[]T) []T {
	res := []T{}
	for _, ts := range s {
		res = append(res, ts...)
	}
	return res
}

func UniqueNumber[T Number](s []T) []T {
	res := []T{}
	unique := map[T]struct{}{}
	for _, ts := range s {
		if _, ok := unique[ts]; !ok {
			res = append(res, ts)
			unique[ts] = struct{}{}
		}
	}
	return res
}

func UniqueMergeNumber[T Number](s ...[]T) []T {
	res := []T{}
	unique := map[T]struct{}{}
	for _, ts := range s {
		for _, t := range ts {
			if _, ok := unique[t]; !ok {
				res = append(res, t)
				unique[t] = struct{}{}
			}
		}
	}
	return res
}

func IntersectNumber[T Number](a []T, s ...[]T) []T {
	intersect := map[T]struct{}{}
	for _, t := range a {
		intersect[t] = struct{}{}
	}
	current := map[T]struct{}{}
	for _, ts := range s {
		for _, t := range ts {
			if _, ok := intersect[t]; ok {
				current[t] = struct{}{}
			}
		}
		intersect = current
		current = map[T]struct{}{}
	}
	res := []T{}
	for t, _ := range intersect {
		res = append(res, t)
	}
	return res
}
