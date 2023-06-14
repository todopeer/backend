package orm

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	TaskStatusNotStarted = iota
	TaskStatusDoing
	TaskStatusDone
	TaskStatusPaused
)

type Task struct {
	ID          int64 `gorm:"primary_key"`
	UserID      *int64
	Name        *string
	Description *string
	Status      *int

	CreatedAt *time.Time
	UpdatedAt *time.Time

	// deprecated
	CompletedAt *time.Time

	DueDate *time.Time
}

func (t *Task) Merge(changes *Task) {
	if changes.UserID != nil {
		t.UserID = changes.UserID
	}
	if changes.Name != nil {
		t.Name = changes.Name
	}
	if changes.Description != nil {
		t.Description = changes.Description
	}
	if changes.Status != nil {
		t.Status = changes.Status
	}
	if changes.CreatedAt != nil {
		t.CreatedAt = changes.CreatedAt
	}
	if changes.UpdatedAt != nil {
		t.UpdatedAt = changes.UpdatedAt
	}
	if changes.CompletedAt != nil {
		t.CompletedAt = changes.CompletedAt
	}
	if changes.DueDate != nil {
		t.DueDate = changes.DueDate
	}
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

// UpdateTask updates an existing task. the `changes` & `user` must be provided
func (t *TaskORM) UpdateTask(current, changes *Task, user *User) error {
	// just in case if changes ID not set, set it
	if changes.ID == 0 {
		changes.ID = current.ID
	}

	now := time.Now()
	runningTaskID := user.RunningTaskID

	return t.db.Transaction(func(tx *gorm.DB) error {
		/*  cases 1: any other status -> doing
		1. if got current running task, update it to paused; update event `end_at` field
		2. set current running task to this one. Create a new Event

			case 2: any other status -> paused/done
		1. if current running task is this one, clear it;
		2. update event `end_at` field, if there's any
		*/
		if changes.Status != nil && *changes.Status == *current.Status {
			changes.Status = nil
		}

		// chagne on status
		if changes.Status != nil {
			if *changes.Status == TaskStatusDoing {
				shouldCreateEvent := true

				if runningTaskID != nil {
					if *runningTaskID != current.ID {
						// update the previous running task to be paused
						if err := tx.Table("tasks").Where("id = ? AND status = ?", *runningTaskID, TaskStatusDoing).Update("status", TaskStatusPaused).Error; err != nil {
							return err
						}

						// create new event for this task
						if err := tx.Table("events").Where("task_id = ? AND end_at IS NULL", *runningTaskID).Update("end_at", &now).Error; err != nil {
							return err
						}
					} else {
						log.Printf("Error: task(id=%d) is running but previous status isn't", current.ID)
						shouldCreateEvent = false
					}
				}

				err := tx.Model(user).Update("running_task_id", current.ID).Error
				if err != nil {
					return err
				}

				if shouldCreateEvent {
					if err := tx.Create(&Event{
						TaskID:  &current.ID,
						StartAt: &now,
					}).Error; err != nil {
						return err
					}
				}
			} else {
				if runningTaskID != nil && *runningTaskID == current.ID {
					if err := tx.Model(user).Update("running_task_id", nil).Error; err != nil {
						return err
					}

					if err := tx.Table("events").Where("task_id = ? AND end_at IS NULL", current.ID).Update("end_at", &now).Error; err != nil {
						return nil
					}
				}
			}
		}

		return tx.Model(current).Update(changes).Error
	})
}

// DeleteTask deletes a task
// if user is passed in, then also check if user.Running task is this task
func (t *TaskORM) DeleteTask(task *Task, user *User) error {
	now := time.Now()
	return t.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(task).Error
		if err != nil {
			return err
		}

		if user != nil && user.RunningTaskID != nil && *user.RunningTaskID == task.ID {
			return tx.Model(&user).Update("running_task_id", nil).Error
		}

		// also mark event as done, as needed
		return tx.Table("events").Where("task_id = ? AND end_at IS NULL", task.ID).Update("end_at", now).Error
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

func GetTasksWithStatus(statuses []int) QueryTaskOptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN (?)", statuses)
	}
}

func GetTasksWithOrder(field string, dir *string) QueryTaskOptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		query := field
		if dir != nil {
			query += " " + *dir
		}
		return db.Order(query, true)
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
