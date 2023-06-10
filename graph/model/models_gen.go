// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type AuthPayload struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type Event struct {
	ID          int64     `json:"id"`
	Task        *Task     `json:"task"`
	Timing      []string  `json:"timing"`
	FullPomo    bool      `json:"fullPomo"`
	TimeCreated time.Time `json:"timeCreated"`
	TimeUpdated time.Time `json:"timeUpdated"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type QueryTaskInput struct {
	Status *TaskStatus `json:"status,omitempty"`
}

type QueryUserTaskResult struct {
	User  *UserPublic `json:"user"`
	Tasks []*Task     `json:"tasks,omitempty"`
	Doing *Task       `json:"doing,omitempty"`
}

type QueryUserTasksInput struct {
	Username string `json:"username"`
}

type Task struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	Events      []*Event   `json:"events,omitempty"`
}

type TaskCreateInput struct {
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

type TaskUpdateInput struct {
	TaskID      int64       `json:"taskId"`
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	DueDate     *time.Time  `json:"dueDate,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
}

type UserPublic struct {
	ID       int64   `json:"id"`
	Username *string `json:"username,omitempty"`
	Name     *string `json:"name,omitempty"`
}

type UserUpdateInput struct {
	ID       int64   `json:"id"`
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
}

type TaskStatus string

const (
	TaskStatusNotStarted TaskStatus = "NOT_STARTED"
	TaskStatusDoing      TaskStatus = "DOING"
	TaskStatusDone       TaskStatus = "DONE"
)

var AllTaskStatus = []TaskStatus{
	TaskStatusNotStarted,
	TaskStatusDoing,
	TaskStatusDone,
}

func (e TaskStatus) IsValid() bool {
	switch e {
	case TaskStatusNotStarted, TaskStatusDoing, TaskStatusDone:
		return true
	}
	return false
}

func (e TaskStatus) String() string {
	return string(e)
}

func (e *TaskStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TaskStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TaskStatus", str)
	}
	return nil
}

func (e TaskStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
