package orm

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/todopeer/backend/util/highorder"
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

	DeletedAt *time.Time

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
func (t *TaskORM) StartTask(task *Task, user *User, eventDesc *string, eventStartAt *time.Time) (*Event, error) {
	/*  cases 1: any other status -> doing
	1. if got current running task, update it to paused; update event `end_at` field
	2. set current running task to this one. Create a new Event
	*/

	if eventStartAt == nil {
		now := time.Now()
		eventStartAt = &now
	}

	var newRunningEventID *int64

	var evt *Event
	err := t.db.Transaction(func(tx *gorm.DB) error {
		return highorder.All(
			// update current running task
			highorder.Branch(user.RunningTaskID != nil,
				// running task is a different one
				highorder.BranchF(func() bool { return *user.RunningTaskID != task.ID }, func() error {
					// stop previous running task & event
					return highorder.All(
						func() error {
							return tx.Table("tasks").Where("id = ? AND status = ?", *user.RunningTaskID, TaskStatusDoing).Update("status", TaskStatusPaused).Error
						},
						func() error {
							return tx.Table("events").Where("task_id = ? AND end_at IS NULL", *user.RunningTaskID).Update("end_at", eventStartAt).Error
						},
					)
				}, func() error {
					// the currently running event
					if user.RunningEventID != nil {
						// update event desc, if provided
						if eventDesc != nil {
							if err := tx.Table("events").Where("id = ?", *user.RunningEventID).Update("description", *eventDesc).Error; err != nil {
								return err
							}
						}
						// and record the running event
						newRunningEventID = user.RunningEventID
					} else {
						log.Printf("unexpected: RunningEventID not set. TaskID: %d", task.ID)
					}
					return nil
				}), nil,
			),

			func() error {
				// task to be doing
				return tx.Table("tasks").Where("id = ?", task.ID).Update("status", TaskStatusDoing).Error
			},

			// update task & event info; only if we need to update Event
			highorder.Branch(newRunningEventID == nil, func() error {
				return highorder.All(
					func() error {
						evt = &Event{
							UserID:      &user.ID,
							TaskID:      &task.ID,
							StartAt:     eventStartAt,
							Description: eventDesc,
						}
						if err := tx.Create(evt).Error; err != nil {
							return err
						}

						newRunningEventID = &evt.ID
						return nil
					},

					func() error {
						userFieldToUpdate := map[string]interface{}{
							"running_task_id":  task.ID,
							"running_event_id": *newRunningEventID,
						}

						return tx.Model(user).Update(userFieldToUpdate).Error
					},
				)
			}, nil),
		)
	})
	return evt, err
}

func (t *TaskORM) UpdateTask(current, changes *Task, user *User) error {
	// just in case if changes ID not set, set it
	if changes.ID == 0 {
		changes.ID = current.ID
	}

	now := time.Now()
	runningTaskID := user.RunningTaskID

	return t.db.Transaction(func(tx *gorm.DB) error {
		/*
				case 2: any other status -> paused/done
			1. if current running task is this one, clear it;
			2. update event `end_at` field, if there's any
		*/
		if changes.Status != nil && *changes.Status == *current.Status {
			changes.Status = nil
		}

		return highorder.All(func() error {
			// make sure we're operating on the doing task
			if changes.Status == nil || runningTaskID == nil || *runningTaskID != current.ID {
				return nil
			}

			// clear the doing status, stop event
			return highorder.All(func() error {
				return tx.Model(user).Update(map[string]any{
					"running_task_id":  nil,
					"running_event_id": nil,
				}).Error
			}, func() error {
				return tx.Table("events").
					Where("task_id = ? AND end_at IS NULL", current.ID).
					Update("end_at", &now).Error
			})
		}, func() error {
			return tx.Model(current).Update(changes).Error
		})
	})
}

// DeleteTask deletes a task
// if user is passed in, then also check if user.Running task is this task
func (t *TaskORM) DeleteTask(task *Task, user *User) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return highorder.All(
			func() error { return tx.Delete(task).Error },
			func() error {
				// running task -- in that case, remove the running event as well
				if user == nil || user.RunningTaskID == nil || *user.RunningTaskID != task.ID {
					return nil
				}

				return highorder.All(
					func() error {
						if user.RunningEventID == nil {
							// actually error case -- when we set the runningTaskID, but not the eventID
							log.Printf("warn: user RunningEvent is nil for runningtask(id=%d)", task.ID)
							return nil
						}

						return tx.Delete(&Event{ID: *user.RunningEventID}).Error
					}, func() error {
						return tx.Model(&user).Update(map[string]interface{}{
							"running_task_id":  nil,
							"running_event_id": nil,
						}).Error
					},
				)
			},
		)
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

type getTaskByIdCfg struct {
	withDeleted bool
}

type GetTasksByIDsOption func(*getTaskByIdCfg)

func GetTasksByIDsOptionWithDeleted(cfg *getTaskByIdCfg) {
	cfg.withDeleted = true
}

func (t *TaskORM) GetTasksByIDs(ids []int64, options ...GetTasksByIDsOption) ([]*Task, error) {
	cfg := &getTaskByIdCfg{}
	for _, option := range options {
		option(cfg)
	}

	var res []*Task
	query := t.db
	if cfg.withDeleted {
		query = query.Unscoped()
	}

	if err := query.Table("tasks").Where("id IN (?)", ids).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
