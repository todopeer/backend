package resolver

import (
	"github.com/flyfy1/diarier/graph/model"
	"github.com/flyfy1/diarier/orm"
	"github.com/flyfy1/diarier/services/auth"
)

func generateToken(user *orm.User) (string, error) {
	return auth.GetTokenFromUser(user)
}
func convertToGraphUserModel(user *orm.User) (*model.User, error) {
	// // TODO: lazy load the Task field instead
	// var task *orm.Task
	// if user.RunningTaskID != nil {
	// 	task, err := r.taskOrm.GetTaskByID(*user.RunningTaskID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return &model.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

// func (r *queryResolver) convertToGraphTaskModel(task *orm.Task) (*model.Task, error) {
// 	var user *orm.User
// 	if task.UserID != 0 {
// 		var err error
// 		user, err = r.userORM.GetUserByID(task.UserID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &model.Task{
// 		ID:            task.ID,
// 		Name:          task.Name,
// 		Description:   task.Description,
// 		Status:        model.Status(task.Status),
// 		DueDate:       task.DueDate,
// 		User:          user,
// 	}, nil
// }
