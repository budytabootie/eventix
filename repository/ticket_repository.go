package repository

import (
    "eventix/entity"
    "gorm.io/gorm"
)

type TicketRepository interface {
    GetAllTickets(page int, size int) ([]entity.Ticket, error)
    GetTicketsByUserID(userID uint) ([]entity.Ticket, error) // Tambahkan metode ini
    GetTicketByID(id uint) (entity.Ticket, error)
    CreateTicket(ticket entity.Ticket) (entity.Ticket, error)
    UpdateTicket(ticket entity.Ticket) (entity.Ticket, error)
    UpdateTicketStatus(id uint, status string) error
    GetSummaryReport() (int64, float64, error)
    GetEventReport(eventID uint) (int64, float64, error)
    SearchTickets(status string) ([]entity.Ticket, error)
    GetPaginatedTickets(page int, size int) ([]entity.Ticket, error)
    IsTicketSold(eventID uint) (bool, error)
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

func (r *ticketRepository) GetEventReport(eventID uint) (int64, float64, error) {
	var totalTickets int64
	var totalRevenue float64

	// Hitung jumlah tiket terjual berdasarkan event
	err := r.db.Model(&entity.Ticket{}).
		Where("event_id = ? AND status = ?", eventID, "purchased").
		Count(&totalTickets).Error
	if err != nil {
		return 0, 0, err
	}

	// Hitung total pendapatan berdasarkan event
	err = r.db.Model(&entity.Ticket{}).
		Where("event_id = ? AND status = ?", eventID, "purchased").
		Select("SUM(price)").Scan(&totalRevenue).Error
	if err != nil {
		return 0, 0, err
	}

	return totalTickets, totalRevenue, nil
}

func (r *ticketRepository) GetSummaryReport() (int64, float64, error) {
	var totalTickets int64
	var totalRevenue float64

	// Hitung jumlah tiket terjual
	err := r.db.Model(&entity.Ticket{}).
		Where("status = ?", "purchased").
		Count(&totalTickets).Error
	if err != nil {
		return 0, 0, err
	}

	// Hitung total pendapatan
	err = r.db.Model(&entity.Ticket{}).
		Where("status = ?", "purchased").
		Select("SUM(price)").Scan(&totalRevenue).Error
	if err != nil {
		return 0, 0, err
	}

	return totalTickets, totalRevenue, nil
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

func (r *ticketRepository) GetTicketsByUserID(userID uint) ([]entity.Ticket, error) {
    var tickets []entity.Ticket
    result := r.db.Where("user_id = ?", userID).Find(&tickets)
    return tickets, result.Error
}

