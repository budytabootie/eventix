package repository

import (
    "eventix/entity"
    "gorm.io/gorm"
)

type TicketRepository interface {
    GetAllTickets(page int, size int) ([]entity.Ticket, error)
    GetTicketsByUserID(userID uint, page int, size int) ([]entity.Ticket, int64, error)
    GetTicketByID(id uint) (entity.Ticket, error)
    CreateTicket(ticket entity.Ticket) (entity.Ticket, error)
    UpdateTicket(ticket entity.Ticket) (entity.Ticket, error)
    UpdateTicketStatus(id uint, status string) error
    GetSummaryReport(page int, size int) ([]entity.Ticket, int64, error) // Update untuk mendukung pagination
    GetEventReport(eventID uint, page int, size int) ([]entity.Ticket, int64, error) // Update untuk mendukung pagination
    SearchTickets(status string) ([]entity.Ticket, error)
    GetPaginatedTickets(page int, size int) ([]entity.Ticket, error)
    IsTicketSold(eventID uint) (bool, error)
    SearchAndFilterTickets(filters map[string]interface{}, page int, size int) ([]entity.Ticket, int64, error) // Tambahkan metode ini
}


type ticketRepository struct {
    db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
    return &ticketRepository{db: db}
}

func (r *ticketRepository) GetAllTickets(page int, size int) ([]entity.Ticket, error) {
    var tickets []entity.Ticket
    offset := (page - 1) * size
    result := r.db.Offset(offset).Limit(size).Find(&tickets)
    return tickets, result.Error
}

func (r *ticketRepository) GetTicketByID(id uint) (entity.Ticket, error) {
    var ticket entity.Ticket
    result := r.db.First(&ticket, id)
    return ticket, result.Error
}

func (r *ticketRepository) CreateTicket(ticket entity.Ticket) (entity.Ticket, error) {
    result := r.db.Create(&ticket)
    return ticket, result.Error
}

func (r *ticketRepository) UpdateTicket(ticket entity.Ticket) (entity.Ticket, error) {
    var existing entity.Ticket
    if err := r.db.First(&existing, ticket.ID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return entity.Ticket{}, gorm.ErrRecordNotFound
        }
        return entity.Ticket{}, err
    }

    // Update hanya pada record yang ditemukan
    if err := r.db.Model(&existing).Updates(ticket).Error; err != nil {
        return entity.Ticket{}, err
    }

    // Mengembalikan record yang diupdate
    return existing, nil
}

func (r *ticketRepository) UpdateTicketStatus(id uint, status string) error {
    result := r.db.Model(&entity.Ticket{}).Where("id = ?", id).Update("status", status)
    return result.Error
}

func (r *ticketRepository) GetEventReport(eventID uint, page int, size int) ([]entity.Ticket, int64, error) {
	var tickets []entity.Ticket
	var totalItems int64

	offset := (page - 1) * size
	query := r.db.Model(&entity.Ticket{}).Where("event_id = ? AND status = ?", eventID, "purchased")

	// Hitung total tiket
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	result := query.Offset(offset).Limit(size).Find(&tickets)
	return tickets, totalItems, result.Error
}

func (r *ticketRepository) GetSummaryReport(page int, size int) ([]entity.Ticket, int64, error) {
	var tickets []entity.Ticket
	var totalItems int64

	offset := (page - 1) * size
	query := r.db.Model(&entity.Ticket{}).Where("status = ?", "purchased")

	// Hitung total tiket
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	result := query.Offset(offset).Limit(size).Find(&tickets)
	return tickets, totalItems, result.Error
}

func (r *ticketRepository) SearchTickets(status string) ([]entity.Ticket, error) {
    var tickets []entity.Ticket
    query := r.db

    if status != "" {
        query = query.Where("status = ?", status)
    }

    result := query.Find(&tickets)
    return tickets, result.Error
}

func (r *ticketRepository) GetPaginatedTickets(page int, size int) ([]entity.Ticket, error) {
    var tickets []entity.Ticket
    offset := (page - 1) * size

    result := r.db.Offset(offset).Limit(size).Find(&tickets)
    return tickets, result.Error
}

func (r *ticketRepository) IsTicketSold(eventID uint) (bool, error) {
    var count int64
    err := r.db.Model(&entity.Ticket{}).Where("event_id = ? AND status = ?", eventID, "sold").Count(&count).Error
    return count > 0, err
}

func (r *ticketRepository) GetTicketsByUserID(userID uint, page int, size int) ([]entity.Ticket, int64, error) {
    var tickets []entity.Ticket
    var totalItems int64

    offset := (page - 1) * size
    query := r.db.Model(&entity.Ticket{}).Where("user_id = ?", userID)

    // Hitung total data sebelum pagination
    if err := query.Count(&totalItems).Error; err != nil {
        return nil, 0, err
    }

    // Pagination
    result := query.Offset(offset).Limit(size).Find(&tickets)
    return tickets, totalItems, result.Error
}


func (r *ticketRepository) SearchAndFilterTickets(filters map[string]interface{}, page int, size int) ([]entity.Ticket, int64, error) {
    var tickets []entity.Ticket
    var totalItems int64

    offset := (page - 1) * size
    query := r.db.Model(&entity.Ticket{})

    // Filter berdasarkan event_id
    if eventID, ok := filters["event_id"].(uint); ok && eventID > 0 {
        query = query.Where("event_id = ?", eventID)
    }

    // Filter berdasarkan status
    if status, ok := filters["status"].(string); ok && status != "" {
        query = query.Where("status = ?", status)
    }

    // Hitung total data sebelum pagination
    if err := query.Count(&totalItems).Error; err != nil {
        return nil, 0, err
    }

    // Ambil data dengan pagination
    result := query.Offset(offset).Limit(size).Find(&tickets)
    return tickets, totalItems, result.Error
}

