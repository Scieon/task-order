package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Task struct {
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"dependencies"`
}

type Tasks struct {
	Tasks []Task `yaml:"tasks"`
}

func main() {
	c, err := readConf("input.yml")
	if err != nil {
		log.Fatal(err)
	}

	taskQueue := make([][]string, 100)
	viewedTasks := make(map[string]int)

	// Current iteration of the task
	taskIteration := 0

	// How many tasks are done
	taskCount := 0

	// while there are still tasks
	for len(c.Tasks) != taskCount {
		for _, task := range c.Tasks {
			if _, ok := viewedTasks[task.Name]; !ok {
				// Immediately add no dependency tasks onto queue
				if len(task.Dependencies) == 0 {
					taskQueue[0] = append(taskQueue[0], task.Name)
					viewedTasks[task.Name] = taskIteration
					taskCount++

					fmt.Println("adding: ", task.Name, " completed in iteration: ", taskIteration)
					continue
				}

				// These handles tasks that have dependencies
				whenCanTaskBeCompleted, canBeDone := canTaskBeDone(viewedTasks, task)

				if canBeDone {
					fmt.Println("adding: ", task.Name, " completed in iteration: ", whenCanTaskBeCompleted)

					taskQueue[whenCanTaskBeCompleted] = append(taskQueue[whenCanTaskBeCompleted], task.Name)
					viewedTasks[task.Name] = whenCanTaskBeCompleted
					taskCount++
				}
			}

		}
		taskIteration++
	}

	// We are done just print out order

	fmt.Println("======================")
	fmt.Println("TASK ORDER")
	fmt.Println("======================")

	for _, tasks := range taskQueue {
		// hacky way of stopping
		if len(tasks) == 0 {
			break
		}

		fmt.Print("[")

		for _, taskName := range tasks {
			fmt.Print(taskName, ", ")
		}
		fmt.Print("]\n")
	}

}

func canTaskBeDone(viewedTasks map[string]int, task Task) (int, bool) {
	highestDep := 0

	// Find if all dependencies has been done
	for _, dep := range task.Dependencies {
		if hasTaskBeenDone(viewedTasks, dep) {
			if highestDep < viewedTasks[dep] {
				highestDep = viewedTasks[dep]
			}
		} else {
			return -1, false
		}
	}

	// Remember to bump by one because we need to first complete the highest depedency task
	return highestDep + 1, true
}

// checks if task has been completed already
func hasTaskBeenDone(viewedTasks map[string]int, taskName string) bool {
	if _, ok := viewedTasks[taskName]; ok {
		return true
	}
	return false
}

func readConf(filename string) (*Tasks, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Tasks{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
