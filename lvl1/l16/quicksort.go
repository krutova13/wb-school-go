package l16

// QuickSort выполняет быструю сортировку массива целых чисел
func QuickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	pivot := arr[0]
	var less, greater []int
	for _, v := range arr[1:] {
		if v <= pivot {
			less = append(less, v)
		} else {
			greater = append(greater, v)
		}
	}
	sortedLess := QuickSort(less)
	withPivot := append(sortedLess, pivot)
	sortedGreater := QuickSort(greater)
	result := append(withPivot, sortedGreater...)
	return result
}
