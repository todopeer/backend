package model

import (
	"github.com/todopeer/backend/orm"
)

func ConvertToGqlPublicUserModel(user *orm.User) *UserPublic {
	return &UserPublic{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		RunningTaskID: user.RunningTaskID,
	}
}

func ConvertToGqlUserModel(user *orm.User) *User {
	return &User{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		Username:      user.Username,
		RunningTaskID: user.RunningTaskID,
	}
}

func ConvertToGqlTaskModel(task *orm.Task) *Task {
	status := AllTaskStatus[*task.Status]

	return &Task{
		ID:          task.ID,
		Name:        *task.Name,
		Description: task.Description,
		Status:      status,
		DueDate:     task.DueDate,
		CreatedAt:   *task.CreatedAt,
		UpdatedAt:   *task.UpdatedAt,
		DeletedAt:   task.DeletedAt,
	}
}

func ConvertToGqlEventModel(e *orm.Event) *Event {
	return &Event{
		ID:          e.ID,
		TaskID:      *e.TaskID,
		StartAt:     *e.StartAt,
		EndAt:       e.EndAt,
		Description: e.Description,
	}
}
