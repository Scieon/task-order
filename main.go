package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	// Keeps track of order of tasks up to 100 iterations
	taskQueue := make([][]string, 100)

	// Keeps track of tasks that have been completed and which iteration it was executed in
	completedTasks := make(map[string]int)

	// Current iteration of the task
	taskIteration := 0

	// How many tasks are done
	taskCount := 0

	// while there are still tasks
	for len(c.Tasks) != taskCount {
		tasksCompletedInCurrentIteration := 0

		// Check if there any task that can be done
		for _, task := range c.Tasks {
			if _, ok := completedTasks[task.Name]; !ok {
				// Immediately add no dependency tasks onto queue
				if len(task.Dependencies) == 0 {
					taskQueue[0] = append(taskQueue[0], task.Name)
					completedTasks[task.Name] = taskIteration
					taskCount++

					fmt.Println( task.Name, " completed in iteration: ", taskIteration)
					tasksCompletedInCurrentIteration++
					continue
				}

				// These handles tasks that have dependencies
				whenCanTaskBeCompleted, canBeDone := canTaskBeDone(completedTasks, task)

				if canBeDone {
					fmt.Println(task.Name, " completed in iteration: ", whenCanTaskBeCompleted)

					taskQueue[whenCanTaskBeCompleted] = append(taskQueue[whenCanTaskBeCompleted], task.Name)
					completedTasks[task.Name] = whenCanTaskBeCompleted
					taskCount++
					tasksCompletedInCurrentIteration++
				}
			}
		}

		// If no tasks are completed in an iteration then we have cyclic dependency
		if tasksCompletedInCurrentIteration == 0 {
			fmt.Println("DETECTED A CYCLIC DEPENDENCY")
			os.Exit(-1)
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

func canTaskBeDone(completedTasks map[string]int, task Task) (int, bool) {
	highestDep := 0

	// Find if all dependencies has been done
	for _, dep := range task.Dependencies {
		if hasTaskBeenDone(completedTasks, dep) {
			if highestDep < completedTasks[dep] {
				highestDep = completedTasks[dep]
			}
		} else {
			return -1, false
		}
	}

	// Remember to bump by one because we need to first complete the highest depedency task
	return highestDep + 1, true
}

// checks if task has been completed already
func hasTaskBeenDone(completedTasks map[string]int, taskName string) bool {
	if _, ok := completedTasks[taskName]; ok {
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
