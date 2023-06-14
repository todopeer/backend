package model

import (
	"github.com/Shopify/hoff"
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

func BuildGqlTaskDetail(task *Task, eventOrm *orm.EventOrm) (*TaskDetail, error) {
	if task == nil {
		return nil, nil
	}

	events, err := eventOrm.GetEventsByTaskID(task.ID)
	if err != nil {
		return nil, err
	}

	return &TaskDetail{
		Task:   task,
		Events: hoff.Map(events, ConvertToGraphqlEvent),
	}, nil
}
