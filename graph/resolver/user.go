package resolver

import (
	"github.com/todopeer/backend/graph/model"
	"github.com/todopeer/backend/orm"
	"github.com/todopeer/backend/services/auth"
)

func generateToken(user *orm.User) (string, error) {
	return auth.GetTokenFromUser(user)
}
func convertToGraphPublicUserModel(user *orm.User) *model.UserPublic {
	return &model.UserPublic{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		RunningTaskID: user.RunningTaskID,
	}
}
func convertToGraphUserModel(user *orm.User) *model.User {
	// // TODO: lazy load the Task field instead
	return &model.User{
		ID:            user.ID,
		Email:         user.Email,
		Username:      user.Username,
		Name:          user.Name,
		RunningTaskID: user.RunningTaskID,
	}
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
