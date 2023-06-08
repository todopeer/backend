package resolver

import (
	"github.com/flyfy1/diarier/graph/model"
	"github.com/flyfy1/diarier/orm"
)

func convertToGraphTaskModel(task *orm.Task) (*model.Task, error) {
	// TODO: lazy load the Event field instead
	var events []*model.Event
	// dbEvents, err := r.EventORM.GetEventsByTaskID(task.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, dbEvent := range dbEvents {
	// 	event, err := r.convertToGraphEventModel(dbEvent)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	events = append(events, event)
	// }

	status := model.AllTaskStatus[*task.Status]

	return &model.Task{
		ID:          task.ID,
		Name:        *task.Name,
		Description: task.Description,
		Status:      status,
		DueDate:     task.DueDate,
		Events:      events,
		CreatedAt:   *task.CreatedAt,
		UpdatedAt:   *task.UpdatedAt,
		CompletedAt: task.CompletedAt,
	}, nil
}
