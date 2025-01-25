package service

import (
	"eventix/repository"
)

type ReportService interface {
	GetSummaryReport(page int, size int) (map[string]interface{}, error) // Update untuk mendukung pagination
	GetEventReport(eventID uint, page int, size int) (map[string]interface{}, error) // Update untuk mendukung pagination
}


type reportService struct {
	ticketRepo repository.TicketRepository
}

func NewReportService(ticketRepo repository.TicketRepository) ReportService {
	return &reportService{
		ticketRepo: ticketRepo,
	}
}

func (s *reportService) GetEventReport(eventID uint, page int, size int) (map[string]interface{}, error) {
	tickets, totalItems, err := s.ticketRepo.GetEventReport(eventID, page, size)
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


func (s *reportService) GetSummaryReport(page int, size int) (map[string]interface{}, error) {
	tickets, totalItems, err := s.ticketRepo.GetSummaryReport(page, size)
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
