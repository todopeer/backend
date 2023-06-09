package orm

import (
	"context"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int64  `gorm:"primary_key"`
	Email         string `gorm:"unique;not null"`
	Name          *string
	Username      *string `gorm:"unique"`
	PasswordHash  string
	RunningTaskID *int64
	SessionID     int32

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserORM struct {
	db *gorm.DB
}

// NewUserORM initializes a new UserORM
func NewUserORM(db *gorm.DB) *UserORM {
	return &UserORM{db: db}
}

func (u *User) HasValidPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// GetUserByEmail retrieves a user by email
func (u *UserORM) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	if err := u.db.Where("email = ?", email).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func (u *UserORM) CreateUser(user *User) error {
	return u.db.Create(user).Error
}

// UpdateUser updates an existing user
func (u *UserORM) UpdateUser(user *User) error {
	return u.db.Save(user).Error
}

// DeleteUser deletes a user
func (u *UserORM) DeleteUser(user *User) error {
	return u.db.Delete(user).Error
}

// GetTasksByUserID retrieves tasks for a specific user
func (t *UserORM) GetUserByID(userID int64) (*User, error) {
	user := &User{}
	if err := t.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (o *UserORM) SetRunningTask(ctx context.Context, user *User, t *Task) error {
	if user.RunningTaskID != nil && *user.RunningTaskID == t.ID {
		// already running this task
		return nil
	}

	return o.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		event := &Event{
			TaskID:  &t.ID,
			StartAt: &now,
		}

		if user.RunningTaskID != nil {
			// re-setting ongoing task, no action needed
			if *user.RunningTaskID == t.ID {
				return nil
			}

			// pause event of previous running task
			err := tx.
				Model(event).
				Where("task_id = ? AND end_at IS NULL", *user.RunningTaskID).
				Update("end_at = ?", now).
				Error

			if err != nil {
				return err
			}
		}

		user.RunningTaskID = &t.ID
		err := tx.Model(user).Update("running_task_id", t.ID).Error
		if err != nil {
			log.Println("db update user error: ", err)
			return err
		}

		err = tx.Model(t).Update("status", ptrInt(TaskStatusDoing)).Error
		if err != nil {
			return err
		}

		// start a new event for this
		return tx.Create(event).Error
	})
}

func (o *UserORM) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	res := &User{}
	err := o.db.Where("username = ?", username).First(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ptrInt(v int) *int {
	return &v
}
