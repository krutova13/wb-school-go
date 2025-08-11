package l8

// SetBit устанавливает i-й бит в числе n в значение bit (0 или 1)
func SetBit(n int64, i uint, bit uint) int64 {
	shift := i - 1
	if bit == 1 {
		return n | (1 << shift)
	}
	return n & ^(1 << shift)
}
