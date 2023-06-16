package model

type User struct {
	ID            int64   `json:"id"`
	Email         string  `json:"email"`
	Name          *string `json:"name,omitempty"`
	Username      *string `json:"username,omitempty"`
	RunningTaskID *int64  `json:"running_task_id,omitempty"`

	// for buffering only
	BufRunningTask *Task
}
