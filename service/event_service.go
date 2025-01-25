package service

import (
	"errors"
	"eventix/entity"
	"eventix/repository"
	"time"
)

type EventService interface {
	GetAllEvents(page int, size int) ([]entity.Event, error)
	GetEventByID(id uint) (entity.Event, error)
	CreateEvent(event entity.Event) (entity.Event, error)
	UpdateEvent(event entity.Event) (entity.Event, error)
	DeleteEvent(id uint) error
	SearchEvents(name string, startDate string, capacity int) ([]entity.Event, error)
}

type eventService struct {
	repo      repository.EventRepository
	ticketRepo repository.TicketRepository
}

func NewEventService(repo repository.EventRepository, ticketRepo repository.TicketRepository) EventService {
	return &eventService{
		repo:      repo,
		ticketRepo: ticketRepo,
	}
}

func (s *eventService) GetAllEvents(page int, size int) ([]entity.Event, error) {
	return s.repo.GetAllEvents(page, size)
}

func (s *eventService) GetEventByID(id uint) (entity.Event, error) {
	return s.repo.GetEventByID(id)
}

func (s *eventService) CreateEvent(event entity.Event) (entity.Event, error) {
	// Validasi kapasitas
	if event.Capacity < 0 {
		return entity.Event{}, errors.New("capacity must be greater than or equal to zero")
	}

	// Validasi harga
	if event.Price < 0 {
		return entity.Event{}, errors.New("price must be greater than or equal to zero")
	}

	// Validasi nama unik
	if err := s.ValidateEventName(event.Name, 0); err != nil {
		return entity.Event{}, err
	}

	// Proses pembuatan event
	return s.repo.CreateEvent(event)
}

func (s *eventService) UpdateEvent(event entity.Event) (entity.Event, error) {
	// Ambil event dari database
	existingEvent, err := s.repo.GetEventByID(event.ID)
	if err != nil {
		return entity.Event{}, errors.New("event not found")
	}

	// Validasi apakah event sudah berlangsung
	if existingEvent.StartDate.Before(time.Now()) {
		return entity.Event{}, errors.New("event cannot be updated because it has already started")
	}

	// Validasi kapasitas
	if event.Capacity < 0 {
		return entity.Event{}, errors.New("capacity must be greater than or equal to zero")
	}

	// Validasi harga
	if event.Price < 0 {
		return entity.Event{}, errors.New("price must be greater than or equal to zero")
	}

	// Proses update event
	return s.repo.UpdateEvent(event)
}

func (s *eventService) DeleteEvent(eventID uint) error {
	// Ambil event dari database
	existingEvent, err := s.repo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	// Validasi apakah event sudah berlangsung
	if existingEvent.StartDate.Before(time.Now()) {
		return errors.New("event cannot be deleted because it has already started")
	}

	// Validasi apakah tiket sudah terjual
	err = s.CanDeleteEvent(eventID)
	if err != nil {
		return err
	}

	// Proses penghapusan event
	return s.repo.DeleteEvent(eventID)
}

func (s *eventService) SearchEvents(name string, startDate string, capacity int) ([]entity.Event, error) {
	return s.repo.SearchEvents(name, startDate, capacity)
}

// Validasi nama unik
func (s *eventService) ValidateEventName(name string, excludeID uint) error {
	isUnique, err := s.repo.IsEventNameUnique(name, excludeID)
	if err != nil {
		return err
	}
	if !isUnique {
		return errors.New("event name must be unique")
	}
	return nil
}

// Validasi tiket yang sudah terjual
func (s *eventService) CanDeleteEvent(eventID uint) error {
	isSold, err := s.ticketRepo.IsTicketSold(eventID)
	if err != nil {
		return err
	}
	if isSold {
		return errors.New("event cannot be deleted because tickets are already sold")
	}
	return nil
}
