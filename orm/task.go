package orm

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	TaskStatusNotStarted = 0
	TaskStatusDoing      = 1
	TaskStatusDone       = 2
)

type Task struct {
	ID          int64 `gorm:"primary_key"`
	UserID      *int64
	Name        *string
	Description *string
	Status      *int

	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	CompletedAt *time.Time

	DueDate *time.Time
}

type TaskORM struct {
	db *gorm.DB
}

// NewTaskORM initializes a new TaskORM
func NewTaskORM(db *gorm.DB) *TaskORM {
	return &TaskORM{db: db}
}

// CreateTask creates a new task
func (t *TaskORM) CreateTask(task *Task) error {
	return t.db.Create(task).Error
}

// UpdateTask updates an existing task
func (t *TaskORM) UpdateTask(task *Task) error {
	return t.db.Update(task).Error
}

// DeleteTask deletes a task
// if user is passed in, then also check if user.Running task is this task
func (t *TaskORM) DeleteTask(task *Task, user *User) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(task).Error
		if err != nil {
			return err
		}

		if user != nil && user.RunningTaskID != nil && *user.RunningTaskID == task.ID {
			return tx.Model(&user).Update("running_task_id", nil).Error
		}
		return nil
	})
}

// GetTaskByID retrieves a task by its ID
func (t *TaskORM) GetTaskByID(id int64) (*Task, error) {
	task := &Task{}
	if err := t.db.First(task, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

type QueryTaskOptionFunc func(db *gorm.DB) *gorm.DB

func GetTasksWithStatus(status int) QueryTaskOptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// GetTasksByUserID retrieves tasks for a specific user
func (t *TaskORM) GetTasksByUserID(userID int64, options ...QueryTaskOptionFunc) ([]*Task, error) {
	var tasks []*Task

	query := t.db.Where("user_id = ?", userID).Order("updated_at DESC")
	for _, option := range options {
		query = option(query)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
