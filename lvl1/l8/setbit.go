package l8

func SetBit(n int64, i uint, bit uint) int64 {
	shift := i - 1
	if bit == 1 {
		return n | (1 << shift)
	}
	return n & ^(1 << shift)
}
