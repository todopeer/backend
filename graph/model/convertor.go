package model

import (
	"github.com/todopeer/backend/orm"
)

func ConvertToGqlPublicUserModel(user *orm.User, taskOrm *orm.TaskORM) *UserPublic {
	return &UserPublic{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		RunningTaskID: user.RunningTaskID,
	}
}

func ConvertToGqlUserModel(user *orm.User, taskOrm *orm.TaskORM) *User {
	eu := NewUser(taskOrm)
	eu.ID = user.ID
	eu.Email = user.Email
	eu.Username = user.Username
	eu.Name = user.Name
	eu.RunningTaskID = user.RunningTaskID

	return eu
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
	}
}

func ConvertToGraphqlEvent(e *orm.Event) *Event {
	return &Event{
		ID:      e.ID,
		StartAt: *e.StartAt,
		EndAt:   e.EndAt,
	}
}
