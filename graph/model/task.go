package model

import (
	"github.com/todopeer/backend/orm"
)

var m = map[TaskStatus]int{
	TaskStatusNotStarted: orm.TaskStatusNotStarted,
	TaskStatusDoing: orm.TaskStatusDoing,
	TaskStatusDone: orm.TaskStatusDone,
}

func TaskStatusToInt(status TaskStatus) int {
	v, found := m[status]
	if found {
		return v
	}

	// cannot find
	return -1
}

func (input *TaskUpdateInput) ChangesAsTask() *orm.Task {
	res := &orm.Task{
		ID:          input.TaskID,
		Name:        input.Name,
		Description: input.Description,
		DueDate:     input.DueDate,
	}

	if input.Status != nil {
		statusInt := TaskStatusToInt(*input.Status)
		res.Status = &statusInt
	}

	return res
}
