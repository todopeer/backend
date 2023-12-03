package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/todopeer/backend/graph"
	"github.com/todopeer/backend/graph/model"
	"github.com/todopeer/backend/services/auth"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.UserRegistrationInput) (*model.User, error) {
	user, err := r.userORM.GetUserByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from DB: %v", err)
	}
	if user != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// create this user
	user = model.ConvertFromRegistrationInput(&input)
	err = r.userORM.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return model.ConvertToGqlUserModel(user), nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	// Fetch the user from the DB
	user, err := r.userORM.GetUserByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from DB: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("no user found with this email")
	}

	// Compare the provided password with the password hash in the DB
	err = user.HasValidPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	// Generate token for the user
	tokenString, err := generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	// Convert the User ORM model to a GraphQL model before returning
	graphUser := model.ConvertToGqlUserModel(user)

	return &model.AuthPayload{
		User:  graphUser,
		Token: tokenString,
	}, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	// In a stateless JWT setup, there's not much you can do here as tokens are generally self-contained.
	// However, if you maintain a blacklist of tokens server-side, you could add the token to that blacklist here.
	// It depends on your setup and needs.

	// For now, we'll just return true indicating the user is logged out.
	user := auth.UserFromContext(ctx)

	// Use your UserORM to fetch the user and increase the SessionID
	dbUser, err := r.userORM.GetUserByEmail(user.Email)
	if err != nil {
		return false, fmt.Errorf("failed to get user from DB: %v", err)
	}

	dbUser.SessionID += 1
	err = r.userORM.UpdateUser(dbUser)
	if err != nil {
		return false, fmt.Errorf("failed to update user in DB: %v", err)
	}

	// The user's JWT token is now invalid because the SessionID stored in the token doesn't match the SessionID in the database.
	return true, nil
}

// UserUpdate is the resolver for the userUpdate field.
func (r *mutationResolver) UserUpdate(ctx context.Context, input model.UserUpdateInput) (*model.User, error) {
	user := auth.UserFromContext(ctx)

	if input.Name != nil {
		user.Name = input.Name
	}

	if input.Username != nil {
		if *input.Username == "" {
			user.Username = &sql.NullString{}
		} else {
			user.Username = &sql.NullString{String: *input.Username, Valid: true}
		}
	}

	if input.Password != nil {
		err := user.SetPassword(*input.Password)
		if err != nil {
			return nil, err
		}
	}

	err := r.userORM.UpdateUser(user)
	return model.ConvertToGqlUserModel(user), err
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	// Get user info from context. The actual implementation depends on how you handle authentication.
	user := auth.UserFromContext(ctx)
	return model.ConvertToGqlUserModel(user), nil
}

// RunningTask is the resolver for the runningTask field.
func (r *userResolver) RunningTask(ctx context.Context, obj *model.User) (*model.Task, error) {
	if obj.RunningTaskID == nil {
		return nil, nil
	}

	if obj.BufRunningTask != nil {
		return obj.BufRunningTask, nil
	}

	t, err := r.taskOrm.GetTaskByID(*obj.RunningTaskID)
	if err != nil {
		return nil, err
	}

	tt := model.ConvertToGqlTaskModel(t)
	if err != nil {
		return nil, err
	}
	obj.BufRunningTask = tt

	return tt, nil
}

// RunningEvent is the resolver for the runningEvent field.
func (r *userResolver) RunningEvent(ctx context.Context, obj *model.User) (*model.Event, error) {
	if obj.RunningEventID == nil {
		return nil, nil
	}

	if obj.BufRunningEvent != nil {
		return obj.BufRunningEvent, nil
	}

	e, err := r.eventOrm.GetEventByID(*obj.RunningEventID)
	if err != nil {
		return nil, err
	}

	ee := model.ConvertToGqlEventModel(e)
	if err != nil {
		return nil, err
	}
	obj.BufRunningEvent = ee

	return ee, nil
}

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
