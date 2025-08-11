package main

import (
	"flag"
	"fmt"
	"os"
	"sortutil/internal"
	"sortutil/internal/sort"
)

func main() {
	setupFlags()

	if err := sort.ProcessSorting(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func setupFlags() {
	flag.IntVar(&internal.GetOpts().Column, "k", 0, "сортировать по столбцу (колонке) N (0 = вся строка)")
	flag.BoolVar(&internal.GetOpts().Numeric, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&internal.GetOpts().Reverse, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&internal.GetOpts().Unique, "u", false, "выводить только уникальные строки")
	flag.BoolVar(&internal.GetOpts().MonthSort, "M", false, "сортировать по названию месяца")
	flag.BoolVar(&internal.GetOpts().TrailingBlanks, "b", false, "игнорировать хвостовые пробелы")
	flag.BoolVar(&internal.GetOpts().CheckSorted, "c", false, "проверить, отсортированы ли данные")
	flag.BoolVar(&internal.GetOpts().HumanReadable, "h", false, "сортировать по размеру")
	flag.Parse()
}
