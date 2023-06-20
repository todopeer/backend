package orm

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	ID     int64
	TaskID *int64
	UserID *int64

	StartAt *time.Time
	EndAt   *time.Time

	Description *string
}

type EventOrm struct {
	db *gorm.DB
}

func NewEventOrm(db *gorm.DB) *EventOrm {
	return &EventOrm{
		db: db,
	}
}

type eventConfig struct {
	limit      *int
	orderAsc   bool
	startAfter *time.Time
}

type EventOptionFunc func(*eventConfig)

func EventQueryOptionWithLimit(limit int) EventOptionFunc {
	return func(ec *eventConfig) {
		ec.limit = &limit
	}
}

func EventQueryOptionWithOrderAsc(toOrderAsc bool) EventOptionFunc {
	return func(ec *eventConfig) {
		ec.orderAsc = toOrderAsc
	}
}

func EventQueryOptionWithStartAfter(startAfter *time.Time) EventOptionFunc {
	return func(ec *eventConfig) {
		ec.startAfter = startAfter
	}
}

type getUserEventsRangeOption struct {
	limit *int32
}
type GetUserEventsRangeOptionFunc func(*getUserEventsRangeOption)

func GetUserEventsRangeWithLimit(limit *int32) GetUserEventsRangeOptionFunc {
	return func(option *getUserEventsRangeOption) {
		option.limit = limit
	}
}

func (e *EventOrm) GetUserEventsRange(userid int64, startTime, endTime time.Time, options ...GetUserEventsRangeOptionFunc) ([]*Event, error) {
	var res []*Event

	cfg := &getUserEventsRangeOption{}
	for _, option := range options {
		option(cfg)
	}

	query := e.db.Table("events").
		Where("user_id = ? AND ( (start_at >= ? AND start_at <= ?) OR (end_at >= ? AND end_at <= ?) OR (start_at <= ? AND end_at IS NULL))",
			userid, startTime, endTime, startTime, endTime, startTime).Order("start_at DESC")

	if cfg.limit != nil {
		query = query.Limit(*cfg.limit)
	}

	if err := query.Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (t *EventOrm) GetEventByID(id int64) (*Event, error) {
	event := &Event{}
	if err := t.db.First(event, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return event, nil
}

// remove an event. Also update the task status if it's running
func (t *EventOrm) DeleteEvent(event *Event) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(event).Error
		if err != nil {
			return err
		}

		if event.EndAt != nil {
			return nil
		}

		// update the task to paused, if it's still doing
		err = tx.Table("tasks").
			Where("id = ? AND status = ?", *event.TaskID, TaskStatusDoing).
			Update("status", TaskStatusPaused).
			Error
		if err != nil {
			return err
		}

		// update user's running task
		return tx.Table("users").
			Where("id = ? AND running_task_id = ?", *event.UserID, *event.TaskID).
			Update("running_task_id", nil).
			Error
	})
}

func (e *EventOrm) GetEventsByTaskID(taskID int64, options ...EventOptionFunc) ([]*Event, error) {
	cfg := &eventConfig{}
	for _, optionFunc := range options {
		optionFunc(cfg)
	}

	var events []*Event
	query := e.db.Model(Event{}).Where("task_id = ?", taskID)
	if cfg.limit != nil {
		query = query.Limit(*cfg.limit)
	}

	if cfg.startAfter != nil {
		query = query.Where("start_at >= ?", cfg.startAfter)
	}

	if cfg.orderAsc {
		query = query.Order("start_at")
	} else {
		query = query.Order("start_at DESC")
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (e *EventOrm) CreateEvent(event *Event) error {
	if event.StartAt == nil {
		now := time.Now()
		event.StartAt = &now
	}

	err := validateStartEndTime(event)
	if err != nil {
		return err
	}

	if event.UserID == nil {
		return errors.New("userid must be defined")
	}

	if event.TaskID == nil {
		return errors.New("taskid must be defined")
	}

	if err := e.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func validateStartEndTime(event *Event) error {
	if event.StartAt == nil {
		return errors.New("start_at must be set")
	}

	if event.StartAt.After(time.Now()) {
		return errors.New("cannot set start_at to be a future value")
	}

	if event.EndAt != nil && event.StartAt.After(*event.EndAt) {
		return errors.New("cannot set start_at to be after end_at")
	}

	return nil
}

// UpdateEvent with optional user field. If passed in, would check:
// 1. if user.runningEvent matches the current event -- if so, pause event if updating event by assigning the eventID
// 2. if updating to a different ID -- if so, update the runningTaskID field on user
func (e *EventOrm) UpdateEvent(event *Event, user *User) error {
	err := validateStartEndTime(event)
	if err != nil {
		return err
	}

	e.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(event).Error; err != nil {
			return err
		}

		if user.RunningEventID != nil && *user.RunningEventID == event.ID {
			if event.EndAt != nil { // event pausing
				// pause the running task
				if user.RunningTaskID != nil {
					err := tx.Table("tasks").Where("id IN (?) AND status = ?", []int64{*user.RunningTaskID, *event.TaskID}, TaskStatusDoing).Update("status", TaskStatusPaused).Error
					if err != nil {
						return err
					}
				}
				// clear the user info
				err = tx.Model(user).Update(map[string]interface{}{
					"running_task_id":  nil,
					"running_event_id": nil,
				}).Error
				if err != nil {
					return err
				}
			} else { // still running
				// event updated to a different taskID
				if user.RunningTaskID != nil && *user.RunningTaskID != *event.TaskID {
					err = tx.Model(user).Update(map[string]interface{}{
						"running_task_id": *event.TaskID,
					}).Error
				}
			}
		}

		return nil
	})

	return nil
}
