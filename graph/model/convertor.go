package model

import (
	"github.com/todopeer/backend/orm"
)

func ConvertToGraphPublicUserModel(user *orm.User, taskOrm *orm.TaskORM) *UserPublic {
	return &UserPublic{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		RunningTaskID: user.RunningTaskID,
	}
}

func ConvertToGraphUserModel(user *orm.User, taskOrm *orm.TaskORM) *User {
	eu := NewUser(taskOrm)
	eu.ID = user.ID
	eu.Email = user.Email
	eu.Username = user.Username
	eu.Name = user.Name
	eu.RunningTaskID = user.RunningTaskID

	return eu
}

func ConvertToGraphTaskModel(task *orm.Task) (*Task, error) {
	status := AllTaskStatus[*task.Status]

	return &Task{
		ID:          task.ID,
		Name:        *task.Name,
		Description: task.Description,
		Status:      status,
		DueDate:     task.DueDate,
		CreatedAt:   *task.CreatedAt,
		UpdatedAt:   *task.UpdatedAt,
	}, nil
}
