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

func (e *EventOrm) GetEventsByTaskID(taskID int) ([]*Event, error) {
	var events []*Event
	if err := e.DB.Where("task_id = ?", taskID).Find(&events).Error; err != nil {
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
