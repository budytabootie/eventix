package service

import (
    "eventix/repository"
)

type ReportService interface {
    GenerateSalesReport() (map[string]interface{}, error)
    GetSummaryReport() (map[string]interface{}, error)
    GetEventReport(eventID uint) (map[string]interface{}, error)
}

type reportService struct {
    ticketRepo repository.TicketRepository
}

// Constructor untuk ReportService dengan dependency injection
func NewReportService(ticketRepo repository.TicketRepository) ReportService {
    return &reportService{ticketRepo: ticketRepo}
}

func (s *reportService) GenerateSalesReport() (map[string]interface{}, error) {
    // Simulasi laporan
    report := map[string]interface{}{
        "total_tickets_sold": 150,
        "total_revenue":      3000,
    }
    return report, nil
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
