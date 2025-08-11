package main

import (
	"cut/internal/cut"
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		fieldsStr = flag.String("f", "", "номера полей для вывода")
		delimiter = flag.String("d", "\t", "разделитель полей")
		separated = flag.Bool("s", false, "только строки, содержащие разделитель")
	)

	flag.Parse()

	if *fieldsStr == "" {
		fmt.Fprintf(os.Stderr, "Ошибка: необходимо указать флаг -f\n")
		flag.Usage()
		os.Exit(1)
	}

	processor, err := cut.NewProcessor(*fieldsStr, *delimiter, *separated)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	if err := processor.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при обработке данных: %v\n", err)
		os.Exit(1)
	}
}
