package model

import (
	"github.com/todopeer/backend/orm"
)

func TaskStatusToInt(status TaskStatus) int {
	for i, defined := range AllTaskStatus {
		if status == defined {
			return i
		}
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
