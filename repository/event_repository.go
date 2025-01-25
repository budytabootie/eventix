package repository

import (
    "eventix/entity"
    "gorm.io/gorm"
)

type EventRepository interface {
    GetAllEvents(page int, size int) ([]entity.Event, error)
    GetEventByID(id uint) (entity.Event, error)
    CreateEvent(event entity.Event) (entity.Event, error)
    UpdateEvent(event entity.Event) (entity.Event, error)
    DeleteEvent(id uint) error
    SearchEvents(name string, startDate string, capacity int) ([]entity.Event, error)
    IsEventNameUnique(name string, excludeID uint) (bool, error)
}

type eventRepository struct {
    db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
    return &eventRepository{db: db}
}

func (r *eventRepository) GetAllEvents(page int, size int) ([]entity.Event, error) {
    var events []entity.Event
    offset := (page - 1) * size
    result := r.db.Offset(offset).Limit(size).Find(&events)
    return events, result.Error
}

func (r *eventRepository) GetEventByID(id uint) (entity.Event, error) {
    var event entity.Event
    result := r.db.First(&event, id)
    return event, result.Error
}

func (r *eventRepository) CreateEvent(event entity.Event) (entity.Event, error) {
    result := r.db.Create(&event)
    return event, result.Error
}

func (r *eventRepository) UpdateEvent(event entity.Event) (entity.Event, error) {
    var existing entity.Event
    if err := r.db.First(&existing, event.ID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return entity.Event{}, gorm.ErrRecordNotFound
        }
        return entity.Event{}, err
    }

    // Update hanya pada record yang ditemukan
    if err := r.db.Model(&existing).Updates(event).Error; err != nil {
        return entity.Event{}, err
    }

    // Mengembalikan record yang diupdate
    return existing, nil
}

func (r *eventRepository) DeleteEvent(id uint) error {
    result := r.db.Delete(&entity.Event{}, id)
    return result.Error
}

func (r *eventRepository) SearchEvents(name string, startDate string, capacity int) ([]entity.Event, error) {
    var events []entity.Event
    query := r.db

    if name != "" {
        query = query.Where("name LIKE ?", "%"+name+"%")
    }
    if startDate != "" {
        query = query.Where("start_date >= ?", startDate)
    }
    if capacity > 0 {
        query = query.Where("capacity >= ?", capacity)
    }

    result := query.Find(&events)
    return events, result.Error
}

func (r *eventRepository) IsEventNameUnique(name string, excludeID uint) (bool, error) {
    var count int64
    query := r.db.Model(&entity.Event{}).Where("name = ?", name)
    if excludeID > 0 {
        query = query.Where("id != ?", excludeID)
    }
    err := query.Count(&count).Error
    return count == 0, err
}
