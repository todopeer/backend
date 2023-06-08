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
		user.RunningTaskID = &t.ID
		err := tx.Model(user).Update("running_task_id", t.ID).Error
		if err != nil {
			log.Println("db update user error: ", err)
			return err
		}

		err = tx.Model(t).Update("status", ptrInt(TaskStatusDoing)).Error
		log.Println("db update task error: ", err)
		return err
	})
}

func ptrInt(v int) *int {
	return &v
}
