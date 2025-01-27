package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

const fileName = "tasks.json"

func main() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nTo-Do List:")
		fmt.Println("1. View Tasks")
		fmt.Println("2. Add Task")
		fmt.Println("3. Mark Task as Complete")
		fmt.Println("4. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("\nYour Tasks:")
			if len(tasks) == 0 {
				fmt.Println("No tasks yet!")
			} else {
				for i, task := range tasks {
					fmt.Printf("%d. %s\n", i+1, task)
				}
			}
		case 2:
			fmt.Print("\nEnter a new task: ")
			scanner.Scan()
			task := scanner.Text()

			// _, err := addTask(tasks, task)
			tasks, err = addTask(tasks, task)
			if err != nil {
				fmt.Println(err)
				continue
			}
			saveTasks(tasks)
			fmt.Println("Task added!")
		case 3:
			fmt.Print("\nEnter the task number to mark as complete: ")
			var taskNum int
			_, err := fmt.Scanln(&taskNum)
			if err != nil {
				fmt.Println("Invalid input. Please enter a valid number.")
				continue
			}

			tasks, err = markTaskComplete(tasks, taskNum)
			if err != nil {
				fmt.Println(err)
				continue
			}
			saveTasks(tasks)
			fmt.Println("Task marked as complete!")
		case 4:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func loadTasks() ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return []string{}, &TaskError{
			Code:    500,
			Message: "Error loading tasks:" + err.Error(),
		}
	}
	defer file.Close()

	var lTasks []string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&lTasks); err != nil {
		return []string{}, &TaskError{
			Code:    500,
			Message: "Error decoding tasks:" + err.Error(),
		}
	}

	return lTasks, nil
}

func saveTasks(Tasks []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return &TaskError{
			Code:    500,
			Message: "Error saving tasks:" + err.Error(),
		}
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(Tasks); err != nil {
		return &TaskError{
			Code:    500,
			Message: "Error encoding tasks:" + err.Error(),
		}
	}

	return nil
}

func markTaskComplete(tasks []string, taskNum int) ([]string, error) {
	if taskNum < 1 || taskNum > len(tasks) {
		return tasks, &TaskError{
			Code:    404,
			Message: "Task number not found. Please try again.",
		}
	}
	tasks = append(tasks[:taskNum-1], tasks[taskNum:]...)
	return tasks, nil
}

func addTask(tasks []string, task string) ([]string, error) {
	if task == "" {
		return tasks, &TaskError{
			Code:    400,
			Message: "Task can't be empty.",
		}
	}
	tasks = append(tasks, task)
	return tasks, nil
}
