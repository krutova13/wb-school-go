package main

import (
	"fmt"
	"math/big"
)

func main() {
	a := big.NewInt(0)
	b := big.NewInt(0)

	a.SetString("478563784956834563", 10)
	b.SetString("21937189204712983423", 10)

	sum := new(big.Int).Add(a, b)
	diff := new(big.Int).Sub(a, b)
	prod := new(big.Int).Mul(a, b)
	quot := new(big.Int).Div(a, b)

	fmt.Printf("a = %s\nb = %s\n", a.String(), b.String())
	fmt.Printf("a + b = %s\n", sum.String())
	fmt.Printf("a - b = %s\n", diff.String())
	fmt.Printf("a * b = %s\n", prod.String())
	fmt.Printf("a / b = %s\n", quot.String())
}
