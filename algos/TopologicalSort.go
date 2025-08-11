package main

import "fmt"

// ScheduleBuild выполняет топологическую сортировку зависимостей проектов
func ScheduleBuild(deps map[string][]string) ([][]string, error) {
	var result [][]string

	for len(deps) > 0 {
		// Находим все проекты без зависимостей в текущем состоянии
		var currentRound []string
		for project, dependencies := range deps {
			if len(dependencies) == 0 {
				currentRound = append(currentRound, project)
			}
		}

		// Добавляем текущий раунд в результат
		result = append(result, currentRound)

		// Удаляем собранные проекты из deps
		for _, project := range currentRound {
			delete(deps, project)
		}

		// Удаляем собранные проекты из зависимостей оставшихся проектов
		for _, project := range currentRound {
			for p, dependencies := range deps {
				newDeps := make([]string, 0, len(dependencies))
				for _, dep := range dependencies {
					if dep != project {
						newDeps = append(newDeps, dep)
					}
				}
				deps[p] = newDeps
			}
		}
	}

	return result, nil
}

func main() {
	deps := map[string][]string{
		"backend":     {"database", "utils"},
		"frontend":    {"utils"},
		"admin-panel": {"backend"},
		"utils":       {},
		"database":    {"utils"},
	}

	schedule, err := ScheduleBuild(deps)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, round := range schedule {
		fmt.Printf("Раунд %d: %v\n", i+1, round)
	}
}
