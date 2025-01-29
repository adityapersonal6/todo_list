package main

import (
	"os"
	"testing"
)

func TestAddTask(t *testing.T) {
	tasks := []string{}
	newTask := "New Task"
	tasks, err := addTask(tasks, newTask)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(tasks) != 1 || tasks[0] != newTask {
		t.Fatalf("Expected task list to contain the new task")
	}
}

func TestAddTask_EmptyTask(t *testing.T) {
	tasks := []string{}
	_, err := addTask(tasks, "")
	if err == nil {
		t.Fatalf("Expected error for empty task, got none")
	}
}

func TestMarkTaskComplete(t *testing.T) {
	tasks := []string{"Task 1", "Task 2", "Task 3"}
	tasks, err := markTaskComplete(tasks, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(tasks) != 2 || tasks[0] != "Task 1" || tasks[1] != "Task 3" {
		t.Fatalf("Expected task list to have the second task removed")
	}
}

func TestMarkTaskComplete_InvalidTaskNum(t *testing.T) {
	tasks := []string{"Task 1", "Task 2", "Task 3"}
	_, err := markTaskComplete(tasks, 4)
	if err == nil {
		t.Fatalf("Expected error for invalid task number, got none")
	}
}

func TestSaveAndLoadTasks(t *testing.T) {
	tasks := []string{"Task 1", "Task 2"}
	err := saveTasks(tasks)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	loadedTasks, err := loadTasks()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(loadedTasks) != len(tasks) {
		t.Fatalf("Expected loaded tasks to match saved tasks")
	}
	for i, task := range tasks {
		if loadedTasks[i] != task {
			t.Fatalf("Expected task %v, got %v", task, loadedTasks[i])
		}
	}

	// Clean up
	os.Remove(fileName)
}
