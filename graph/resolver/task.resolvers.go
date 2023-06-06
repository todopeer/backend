package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/flyfy1/diarier/graph/model"
	"github.com/flyfy1/diarier/orm"
	"github.com/flyfy1/diarier/services/auth"
)

// TaskCreate is the resolver for the taskCreate field.
func (r *mutationResolver) TaskCreate(ctx context.Context, input model.TaskCreateInput) (*model.Task, error) {
	user := auth.UserFromContext(ctx)

	status := orm.TaskStatusNotStarted
	dbTask := &orm.Task{
		Name:        &input.Name,
		UserID:      &user.ID,
		Description: input.Description,
		Status:      &status,
		DueDate:     input.DueDate,
	}

	if err := r.taskOrm.CreateTask(dbTask); err != nil {
		return nil, err
	}

	return convertToGraphTaskModel(dbTask)
}

// TaskUpdate is the resolver for the taskUpdate field.
func (r *mutationResolver) TaskUpdate(ctx context.Context, input model.TaskUpdateInput) (*model.Task, error) {
	user := auth.UserFromContext(ctx)
	task, err := r.taskOrm.GetTaskByID(input.TaskID)
	if err != nil {
		return nil, err
	}

	if *task.UserID != user.ID {
		return nil, errors.New("not authorized to update this task")
	}

	changes := input.ChangesAsTask()
	now := time.Now()
	if changes.Status != nil && *changes.Status == orm.TaskStatusDone && *task.Status != orm.TaskStatusDone {
		// status changed to done
		changes.CompletedAt = &now
	}
	// TODO: consider clear completed_at, if undone a task

	if err := r.taskOrm.UpdateTask(changes); err != nil {
		return nil, err
	}

	return convertToGraphTaskModel(task)
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context, date time.Time) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented: Events - events"))
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context, input model.QueryTaskInput) ([]*model.Task, error) {
	user := auth.UserFromContext(ctx)

	var options []orm.QueryTaskOptionFunc
	if input.Status != nil {
		options = append(options, orm.GetTasksWithStatus(model.TaskStatusToInt(*input.Status)))
	}

	dbTasks, err := r.taskOrm.GetTasksByUserID(user.ID, options...)
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	for _, dbTask := range dbTasks {
		task, err := convertToGraphTaskModel(dbTask)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
