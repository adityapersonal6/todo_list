package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/rivo/tview"
)

const fileName = "tasks.json"

func main() {

	// Create a new application.
	app := tview.NewApplication()

	// Load the tasks.
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println(err)
	}

	// Create a TextView that will display the tasks.
	textList := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	//  Set border and title of the TextView.
	textList.SetBorder(true).SetTitle("To-Do List")

	// refreshes the tasks display whenever there are changes
	refreshTasks := func() {
		textList.Clear()

		if len(tasks) == 0 {
			fmt.Fprintf(textList, "No tasks yet!")
		} else {
			for i, task := range tasks {
				fmt.Fprintf(textList, "[%d] %s\n", i+1, task)
			}
		}
	}

	// Create a form for adding tasks.
	taskNameInput := tview.NewInputField().SetLabel("Task: ")
	// Create a form for marking tasks as complete.
	taskMarkAsCompleteInput := tview.NewInputField().SetLabel("Task ID to Complete: ")

	// Create a form and add the input fields and buttons to it.
	form := tview.NewForm().
		AddFormItem(taskNameInput).
		AddFormItem(taskMarkAsCompleteInput).
		AddButton("Add Task", func() {
			tasks, err = addTask(tasks, taskNameInput.GetText())
			if err != nil {
				fmt.Fprintln(textList, "\n\nWarning!\n", err)
				return
			}
			saveTasks(tasks)
			refreshTasks()
		}).
		AddButton("Mark as Complete", func() {
			taskNumStr := taskMarkAsCompleteInput.GetText()
			if taskNumStr == "" {
				fmt.Fprintln(textList, "\n\nWarning!\nPlease enter a task ID to mark as complete.")
				return
			}
			taskNum, err := strconv.Atoi(taskNumStr)
			if err != nil {
				fmt.Fprintln(textList, "\n\nWarning!\nInvalid task number.")
				return
			}
			tasks, err = markTaskComplete(tasks, taskNum)
			if err != nil {
				fmt.Fprintln(textList, "\n\nWarning!\n", err)
				return
			}
			saveTasks(tasks)
			refreshTasks()
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	// Set the form's border and title.
	form.SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	// Create a flex layout and add the TextView and Form to it.
	flex := tview.NewFlex().
		AddItem(textList, 0, 1, false).
		AddItem(form, 0, 1, true)

	// Refresh the tasks display.
	refreshTasks()

	// Set the root of the application to the flex layout.
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
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
