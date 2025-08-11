package sort

import (
	"fmt"
	"os"
	"sortutil/internal"
)

// ProcessSorting выполняет основную логику сортировки.
// Читает входные данные, проверяет опции сортировки,
// выполняет сортировку и выводит результат
func ProcessSorting() error {
	lines, err := internal.ReadInput()
	if err != nil {
		return fmt.Errorf("ошибка чтения: %v", err)
	}

	if internal.GetOpts().CheckSorted {
		sortableLines := internal.PrepareLines(lines)
		if sortableLines.IsSorted() {
			fmt.Println("Данные уже отсортированы")
			return nil
		}
		fmt.Println("Данные не отсортированы")
		os.Exit(1)
	}

	sortableLines := internal.PrepareLines(lines)
	sortableLines.Sort()
	internal.OutputLines(sortableLines)

	return nil
}
