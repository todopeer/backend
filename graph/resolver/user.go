package resolver

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/flyfy1/diarier/graph/model"
	"github.com/flyfy1/diarier/orm"
)

const (
	jwtKey            = "top-secret-key" // Use a strong secret key for your JWT tokens
	jwtExpireDuration = 72 * time.Hour   // Tokens will expire after 72 hours
)

func generateToken(user *orm.User) (string, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(jwtExpireDuration).Unix(),
		Issuer:    "diarier",
		Subject:   fmt.Sprintf("%d", user.ID),
		Id:        fmt.Sprintf("%d", user.SessionID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
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
