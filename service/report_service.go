package service

import (
	"eventix/repository"
)

type ReportService interface {
	GetSummaryReport() (map[string]interface{}, error)
	GetEventReport(eventID uint) (map[string]interface{}, error)
}

type reportService struct {
	ticketRepo repository.TicketRepository
}

func NewReportService(ticketRepo repository.TicketRepository) ReportService {
	return &reportService{
		ticketRepo: ticketRepo,
	}
}

func (s *reportService) GetEventReport(eventID uint) (map[string]interface{}, error) {
	totalTickets, totalRevenue, err := s.ticketRepo.GetEventReport(eventID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_tickets_sold": totalTickets,
		"total_revenue":      totalRevenue,
	}, nil
}

func (s *reportService) GetSummaryReport() (map[string]interface{}, error) {
	totalTickets, totalRevenue, err := s.ticketRepo.GetSummaryReport()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_tickets_sold": totalTickets,
		"total_revenue":      totalRevenue,
	}, nil
}
