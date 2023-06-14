package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/todopeer/backend/graph/model"
	"github.com/todopeer/backend/orm"
	"github.com/todopeer/backend/services/auth"
	"github.com/todopeer/backend/util"
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

	return model.ConvertToGraphTaskModel(dbTask)
}

// TaskUpdate is the resolver for the taskUpdate field.
func (r *mutationResolver) TaskUpdate(ctx context.Context, id int64, input model.TaskUpdateInput) (*model.Task, error) {
	user := auth.UserFromContext(ctx)
	task, err := r.taskOrm.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	if *task.UserID != user.ID {
		return nil, errors.New("not authorized to update this task")
	}

	changes := input.ChangesAsTask()
	if err := r.taskOrm.UpdateTask(task, changes, user); err != nil {
		return nil, err
	}

	return model.ConvertToGraphTaskModel(task)
}

// TaskRemove is the resolver for the taskRemove field.
func (r *mutationResolver) TaskRemove(ctx context.Context, id int64) (*model.Task, error) {
	user := auth.UserFromContext(ctx)
	task, err := r.taskOrm.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if *task.UserID != user.ID {
		return nil, errors.New("not authorized to update this task")
	}

	err = r.taskOrm.DeleteTask(task, user)
	if err != nil {
		return nil, err
	}

	return model.ConvertToGraphTaskModel(task)
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context, input model.QueryTaskInput) ([]*model.Task, error) {
	user := auth.UserFromContext(ctx)

	var options []orm.QueryTaskOptionFunc
	if len(input.Status) > 0 {
		statuses := util.Map(input.Status, model.TaskStatusToInt)
		options = append(options, orm.GetTasksWithStatus(statuses))
	}

	if input.OrderBy != nil {
		options = append(options, orm.GetTasksWithOrder(string(input.OrderBy.Field), (*string)(input.OrderBy.Direction)))
	}

	dbTasks, err := r.taskOrm.GetTasksByUserID(user.ID, options...)
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	for _, dbTask := range dbTasks {
		task, err := model.ConvertToGraphTaskModel(dbTask)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UserTasks is the resolver for the userTasks field.
func (r *queryResolver) UserTasks(ctx context.Context, username string) (*model.QueryUserTaskResult, error) {
	user, err := r.userORM.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// load task under this user
	tasks, err := r.taskOrm.GetTasksByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	taskResp, err := util.MapWithError(tasks, model.ConvertToGraphTaskModel)
	if err != nil {
		return nil, err
	}

	res := &model.QueryUserTaskResult{
		User:  model.ConvertToGraphPublicUserModel(user, r.taskOrm),
		Tasks: taskResp,
	}

	if user.RunningTaskID != nil {
		idx := util.FindBy(taskResp, func(t *model.Task) bool {
			return t.ID == *user.RunningTaskID
		})
		if idx >= 0 {
			res.Doing = taskResp[idx]
		} else {
			log.Printf("RunningTaskID(=%d) exist but not in this user's task list", *user.RunningTaskID)
		}
	}
	return res, nil
}
