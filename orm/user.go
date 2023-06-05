package orm

import (
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
	if err := t.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
