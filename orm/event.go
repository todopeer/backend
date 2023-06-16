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
	DB *gorm.DB
}

func NewEventOrm(db *gorm.DB) *EventOrm {
	return &EventOrm{
		DB: db,
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

func (e *EventOrm) GetUserEventsByDay(userid int64, dayStart time.Time) ([]*Event, error) {
	startTime := dayStart
	endTime := startTime.Add(time.Hour * 24)
	var res []*Event

	if err := e.DB.Table("events").Where("user_id = ? AND start_at >= ? AND end_at < ?", userid, startTime, endTime).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (e *EventOrm) GetEventsByTaskID(taskID int64, options ...EventOptionFunc) ([]*Event, error) {
	cfg := &eventConfig{}
	for _, optionFunc := range options {
		optionFunc(cfg)
	}

	var events []*Event
	query := e.DB.Model(Event{}).Where("task_id = ?", taskID)
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

	if err := e.DB.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (e *EventOrm) UpdateEvent(event *Event) error {
	if err := e.DB.Save(event).Error; err != nil {
		return err
	}
	return nil
}
