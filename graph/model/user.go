package model

import (
	"github.com/todopeer/backend/orm"
)

type User struct {
	taskOrm       *orm.TaskORM
	ID            int64   `json:"id"`
	Email         string  `json:"email"`
	Name          *string `json:"name,omitempty"`
	Username      *string `json:"username,omitempty"`
	RunningTaskID *int64  `json:"running_task_id,omitempty"`

	// for buffering only
	BufRunningTask *Task
}

func NewUser(orm *orm.TaskORM) *User {
	return &User{taskOrm: orm}
}
