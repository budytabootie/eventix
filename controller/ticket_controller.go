package controller

import (
	"eventix/entity"
	"eventix/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	service service.TicketService
}

func NewTicketController(ticketService service.TicketService) *TicketController {
	return &TicketController{
		service: ticketService,
	}
}


// GetTickets godoc
// @Summary Get user tickets
// @Description Retrieve tickets owned by the logged-in user
// @Tags Tickets
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tickets [get]
func (ctrl *TicketController) GetTickets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	tickets, err := ctrl.service.GetTicketsByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve tickets", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Tickets retrieved successfully", "data": tickets})
}


// CreateTicket godoc
// @Summary Purchase a ticket
// @Description User can purchase tickets for an event
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticket body entity.Ticket true "Ticket details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tickets [post]
func (ctrl *TicketController) CreateTicket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	var ticket entity.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ticket data", "data": nil})
		return
	}

	// Assign user_id dari sesi login
	ticket.UserID = userID.(uint)

	createdTicket, err := ctrl.service.CreateTicket(ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Ticket created successfully", "data": createdTicket})
}


func (ctrl *TicketController) CancelTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ticket ID", "data": nil})
		return
	}

	if err := ctrl.service.CancelTicket(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Ticket cancelled successfully", "data": nil})
}

func (ctrl *TicketController) GetTicketByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ticket ID", "data": nil})
		return
	}
	ticket, err := ctrl.service.GetTicketByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve ticket", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Ticket retrieved successfully", "data": ticket})
}

func (ctrl *TicketController) UpdateTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ticket ID", "data": nil})
		return
	}

	var ticket entity.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ticket data", "data": nil})
		return
	}
	ticket.ID = uint(id)

	updatedTicket, err := ctrl.service.UpdateTicket(ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Ticket updated successfully", "data": updatedTicket})
}

