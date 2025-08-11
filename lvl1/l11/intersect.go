package l11

// Intersect находит пересечение двух массивов целых чисел
func Intersect(arr1, arr2 []int) []int {
	set := make(map[int]struct{})
	var result []int

	for _, v := range arr1 {
		set[v] = struct{}{}
	}

	for _, v := range arr2 {
		if _, ok := set[v]; ok {
			result = append(result, v)
			delete(set, v)
		}
	}
	return result
}
