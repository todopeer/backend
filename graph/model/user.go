package model

type User struct {
	ID    int64   `json:"id"`
	Email string  `json:"email"`
	Name  *string `json:"name,omitempty"`
}

func (u *User) RunningTask() *Task {
	// TODO: load dynamically
	return nil
}