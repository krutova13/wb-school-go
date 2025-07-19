package main

import (
	"fmt"
	"wbschoolgo/lvl1/l21"
)

func main() {
	legacy := &l21.LegacyReader{}
	adapter := l21.NewLegacyReaderAdapter(legacy)
	fmt.Println(adapter.ReadData())
}
