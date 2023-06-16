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

func (e *EventOrm) GetUserEventsRange(userid int64, startTime, endTime time.Time) ([]*Event, error) {
	var res []*Event

	if err := e.db.Table("events").Where("user_id = ? AND ( (start_at >= ? AND start_at <= ?) OR (end_at >= ? AND end_at <= ?) OR (start_at <= ? AND end_at IS NULL))", userid,
		startTime, endTime, startTime, endTime, startTime).Find(&res).Error; err != nil {

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

func (e *EventOrm) UpdateEvent(event *Event) error {
	if event.StartAt == nil {
		return errors.New("start_at must be set")
	}

	if event.StartAt.After(time.Now()) {
		return errors.New("cannot set start_at to be a future value")
	}

	if event.EndAt != nil && event.StartAt.After(*event.EndAt) {
		return errors.New("cannot set start_at to be after end_at")
	}

	if err := e.db.Save(event).Error; err != nil {
		return err
	}
	return nil
}
