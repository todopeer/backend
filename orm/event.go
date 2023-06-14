package orm

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	ID     int64
	TaskID *int64

	StartAt *time.Time
	EndAt   *time.Time
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
