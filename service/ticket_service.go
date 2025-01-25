package service

import (
	"errors"
	"eventix/entity"
	"eventix/repository"
)

type TicketService interface {
    GetAllTickets(page int, size int) ([]entity.Ticket, error)
    GetTicketsByUserID(userID uint, page int, size int) ([]entity.Ticket, int64, error) // Perbaikan di sini
    GetTicketByID(id uint) (entity.Ticket, error)
    CreateTicket(ticket entity.Ticket) (entity.Ticket, error)
    UpdateTicket(ticket entity.Ticket) (entity.Ticket, error)
    CancelTicket(ticketID uint) error
    UpdateTicketStatus(id uint, status string) error
    SearchTickets(status string) ([]entity.Ticket, error)
    GetPaginatedTickets(page int, size int) ([]entity.Ticket, error)
	SearchAndFilterTickets(filters map[string]interface{}, page int, size int) (map[string]interface{}, error)
}



type ticketService struct {
	repo      repository.TicketRepository
	eventRepo repository.EventRepository
}

func NewTicketService(repo repository.TicketRepository, eventRepo repository.EventRepository) TicketService {
	return &ticketService{repo: repo, eventRepo: eventRepo}
}

func (s *ticketService) GetAllTickets(page int, size int) ([]entity.Ticket, error) {
	return s.repo.GetAllTickets(page, size)
}

func (s *ticketService) GetTicketsByUserID(userID uint, page int, size int) ([]entity.Ticket, int64, error) {
    return s.repo.GetTicketsByUserID(userID, page, size)
}


func (s *ticketService) GetTicketByID(id uint) (entity.Ticket, error) {
	return s.repo.GetTicketByID(id)
}

func (s *ticketService) CreateTicket(ticket entity.Ticket) (entity.Ticket, error) {
	// Ambil data event terkait
	event, err := s.eventRepo.GetEventByID(ticket.EventID)
	if err != nil {
		return entity.Ticket{}, errors.New("event not found")
	}

	// Validasi kapasitas
	if ticket.Quantity > event.Capacity {
		return entity.Ticket{}, errors.New("quantity exceeds event capacity")
	}

	// Hitung harga tiket berdasarkan harga event
	ticket.Price = float64(ticket.Quantity) * event.Price
	ticket.Status = "purchased"

	// Kurangi kapasitas event
	event.Capacity -= ticket.Quantity
	_, err = s.eventRepo.UpdateEvent(event)
	if err != nil {
		return entity.Ticket{}, errors.New("failed to update event capacity")
	}

	// Buat tiket
	return s.repo.CreateTicket(ticket)
}


func (s *ticketService) CancelTicket(ticketID uint) error {
	// Ambil data tiket terkait
	ticket, err := s.repo.GetTicketByID(ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	// Validasi status tiket
	if ticket.Status != "purchased" {
		return errors.New("ticket cannot be cancelled")
	}

	// Ambil data event terkait
	event, err := s.eventRepo.GetEventByID(ticket.EventID)
	if err != nil {
		return errors.New("event not found")
	}

	// Update kapasitas event
	event.Capacity += ticket.Quantity
	_, err = s.eventRepo.UpdateEvent(event)
	if err != nil {
		return errors.New("failed to update event capacity")
	}

	// Update status tiket
	return s.repo.UpdateTicketStatus(ticket.ID, "cancelled")
}

func (s *ticketService) UpdateTicketStatus(id uint, status string) error {
	_, err := s.repo.GetTicketByID(id)
	if err != nil {
		return err
	}

	return s.repo.UpdateTicketStatus(id, status)
}

func (s *ticketService) SearchTickets(status string) ([]entity.Ticket, error) {
	return s.repo.SearchTickets(status)
}

func (s *ticketService) GetPaginatedTickets(page int, size int) ([]entity.Ticket, error) {
	return s.repo.GetPaginatedTickets(page, size)
}

func (s *ticketService) UpdateTicket(ticket entity.Ticket) (entity.Ticket, error) {
    // Ambil data tiket terkait
    existingTicket, err := s.repo.GetTicketByID(ticket.ID)
    if err != nil {
        return entity.Ticket{}, errors.New("ticket not found")
    }

    // Validasi status tiket
    if existingTicket.Status != "purchased" {
        return entity.Ticket{}, errors.New("ticket cannot be updated unless it is in 'purchased' status")
    }

    // Perbarui data tiket
    updatedTicket, err := s.repo.UpdateTicket(ticket)
    if err != nil {
        return entity.Ticket{}, err
    }

    return updatedTicket, nil
}

func (s *ticketService) SearchAndFilterTickets(filters map[string]interface{}, page int, size int) (map[string]interface{}, error) {
    tickets, totalItems, err := s.repo.SearchAndFilterTickets(filters, page, size)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "tickets":      tickets,
        "total_items":  totalItems,
        "current_page": page,
        "page_size":    size,
    }, nil
}
