package model

import (
	"github.com/todopeer/backend/orm"
)

func ConvertToGqlPublicUserModel(user *orm.User) *UserPublic {
	if user == nil {
		return nil
	}

	return &UserPublic{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		RunningTaskID: user.RunningTaskID,
	}
}

func ConvertToGqlUserModel(user *orm.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		Username:       user.Username,
		RunningTaskID:  user.RunningTaskID,
		RunningEventID: user.RunningEventID,
	}
}

func ConvertToGqlTaskModel(task *orm.Task) *Task {
	if task == nil {
		return nil
	}

	status := AllTaskStatus[*task.Status]
	r := &Task{
		ID:          task.ID,
		Name:        *task.Name,
		Description: task.Description,
		Status:      status,
		DueDate:     task.DueDate,
		CreatedAt:   *task.CreatedAt,
		UpdatedAt:   *task.UpdatedAt,
	}
	if task.DeletedAt != nil && task.DeletedAt.Valid {
		r.DeletedAt = &task.DeletedAt.Time
	}

	return r
}

func ConvertToGqlEventModel(e *orm.Event) *Event {
	if e == nil {
		return nil
	}

	return &Event{
		ID:          e.ID,
		TaskID:      *e.TaskID,
		StartAt:     *e.StartAt,
		EndAt:       e.EndAt,
		Description: e.Description,
	}
}
