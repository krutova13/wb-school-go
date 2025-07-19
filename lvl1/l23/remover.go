package l23

import "fmt"

func RemoveElementByIndex(slice []int, index int) ([]int, error) {
	if index < 0 || index >= len(slice) {
		return slice, fmt.Errorf("индекс %d вне диапазона слайса [0,%d]", index, len(slice))
	}
	copy(slice[index:], slice[index+1:])
	slice[len(slice)-1] = 0
	return slice[:len(slice)-1], nil
}
