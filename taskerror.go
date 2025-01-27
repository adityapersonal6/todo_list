package main

import "fmt"

type TaskError struct {
    Code    int
    Message string
}

func (e *TaskError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}