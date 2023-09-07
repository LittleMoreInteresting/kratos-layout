package slicehelper

import (
	"fmt"
	"testing"
)

func TestInSlice(t *testing.T) {

	fmt.Println(InSlice(1, []int{}) == false)
	fmt.Println(InSlice(1, []int{2, 3}) == false)
	fmt.Println(InSlice(4, []int{9, 8, 4, 2, 3}) == true)
	fmt.Println(InSlice(1, []int64{2, 3}) == false)
	fmt.Println(InSlice(4, []int64{9, 8, 4, 2, 3}) == true)
	fmt.Println(InSlice(10, []int64{9, 8, 4, 2, 3}) == false)
	fmt.Println(InSlice(9, []int64{9, 8, 4, 2, 3}) == true)
	fmt.Println(InSlice(9, []int32{9, 8, 4, 2, 3}) == true)
	fmt.Println(InSlice(9, []float64{9, 8, 4, 2, 3}) == true)
	fmt.Println(InSlice(2, []int{9, 8, 4, 2, 3}) == true)
	fmt.Println(InSlice(2, []int64{9, 8, 4, 2, 3}) == true)
	fmt.Println(InSlice(2, []float64{9, 8, 4, 2, 3}) == true)
}

func TestMergeSlice(t *testing.T) {
	fmt.Println(SimpleMergeNumber([]int{}))
	fmt.Println(SimpleMergeNumber([]int{1, 2}, []int{3, 4}, []int{5, 6, 7}, []int{8}))
	fmt.Println(SimpleMergeNumber([]int32{1, 2}, []int32{3, 4}, []int32{5, 6, 7}, []int32{8}))
	fmt.Println(SimpleMergeNumber([]int64{1, 2}, []int64{3, 4}, []int64{5, 6, 7}, []int64{8}))
}

func TestUniqueMergeNumber(t *testing.T) {
	fmt.Println(UniqueMergeNumber([]int{}))
	fmt.Println(UniqueMergeNumber([]int{1, 2}, []int{3, 4}, []int{1, 2, 3, 5, 6, 7}, []int{8}))
	fmt.Println(UniqueMergeNumber([]int32{1, 2}, []int32{3, 4}, []int32{1, 2, 3, 5, 6, 7}, []int32{8}))
	fmt.Println(UniqueMergeNumber([]int64{1, 2}, []int64{3, 4}, []int64{1, 2, 3, 5, 6, 7}, []int64{8}))
}

func TestUniqueNumber(t *testing.T) {
	fmt.Println(UniqueNumber([]int{}))
	fmt.Println(UniqueNumber([]int{1, 2, 3, 5, 6, 7, 2, 3, 5}))
	fmt.Println(UniqueNumber([]int32{1, 2, 3, 5, 6, 7, 2, 3, 5}))
	fmt.Println(UniqueNumber([]int64{1, 2, 3, 5, 6, 7, 2, 3, 5}))
}

func TestIntersectNumber(t *testing.T) {
	fmt.Println(IntersectNumber([]int{1, 2}, []int{1, 2, 9}, []int{1, 2, 3, 5, 6, 7}))
	fmt.Println(IntersectNumber([]int{1, 2}, []int{3, 3, 4}, []int{1, 2, 3, 5, 6, 7}))
	fmt.Println(IntersectNumber([]int{1, 2}, []int{2, 3, 3, 4}, []int{1, 2, 3, 5, 6, 7}))
}
