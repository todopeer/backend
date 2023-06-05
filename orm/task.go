package orm

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	ID          int64 `gorm:"primary_key"`
	UserID      int64
	Name        string
	Description *string
	Status      int

	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	return t.db.Save(task).Error
}

// DeleteTask deletes a task
func (t *TaskORM) DeleteTask(task *Task) error {
	return t.db.Delete(task).Error
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

// GetTasksByUserID retrieves tasks for a specific user
func (t *TaskORM) GetTasksByUserID(userID int64) ([]*Task, error) {
	var tasks []*Task
	if err := t.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
