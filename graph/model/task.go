package model

import (
	"github.com/todopeer/backend/orm"
)

func TaskStatusToInt(status TaskStatus) int {
	for i, s := range AllTaskStatus {
		if s == status {
			return i
		}
	}

	return -1
}

func (input *TaskUpdateInput) ChangesAsTask() *orm.Task {
	res := &orm.Task{
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
