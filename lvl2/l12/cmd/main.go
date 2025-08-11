package main

import (
	"flag"
	"fmt"
	"greputil/internal"
	"greputil/internal/grep"
	"os"
)

func main() {
	setupFlags()

	if err := grep.ProcessGrepping(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func setupFlags() {
	flag.IntVar(&internal.GetOpts().After, "A", 0, "после каждой найденной строки дополнительно вывести N строк после неё")
	flag.IntVar(&internal.GetOpts().Before, "B", 0, "вывести N строк до каждой найденной строки")
	flag.IntVar(&internal.GetOpts().Circle, "C", 0, "вывести N строк контекста вокруг найденной строки")
	flag.BoolVar(&internal.GetOpts().Count, "c", false, "выводить только то количество строк, что совпадающих с шаблоном")
	flag.BoolVar(&internal.GetOpts().Ignore, "i", false, "игнорировать регистр")
	flag.BoolVar(&internal.GetOpts().Invert, "v", false, "инвертировать фильтр: выводить строки, не содержащие шаблон")
	flag.BoolVar(&internal.GetOpts().Fix, "F", false, "воспринимать шаблон как фиксированную строку, а не регулярное выражение")
	flag.BoolVar(&internal.GetOpts().Number, "n", false, "выводить номер строки перед каждой найденной строкой")
	flag.Parse()
}
