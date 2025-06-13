package controllers

import (
	"net/http"
	"slices"

	"github.com/Faith-Kiv/Ticketing-Backend/models"
	"github.com/Faith-Kiv/Ticketing-Backend/services"
	"github.com/gin-gonic/gin"
)

func GetTickets(ctx *gin.Context) {
	tickets := []models.Tickets{}

	// Simulate fetching tickets from a database or service
	if err := ctx.ShouldBindQuery(&tickets); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
}

func CreateTicket(ctx *gin.Context) {
	roles := ctx.GetStringSlice("roles")
	ticket := models.Tickets{}

	if err := ctx.ShouldBindJSON(&ticket); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	isCustomerSupport := slices.Contains(roles, "customer_support")
	if !isCustomerSupport {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: You do not have permission to create a ticket."})
		return
	}

	if !services.ValidatePhoneNumber(ticket.CustomerPhoneNumber) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}

}

func GetTicket(ctx *gin.Context) {
	ticketID := ctx.Param("id")

	if ticketID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

}

func UpdateTicket(ctx *gin.Context) {
	ticketID := ctx.Param("id")
	if ticketID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	ticket := models.Tickets{}
	if err := ctx.ShouldBindJSON(&ticket); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Simulate updating the ticket in a database or service
	ctx.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully", "ticket_id": ticketID})
}

func DeleteTicket(ctx *gin.Context) {
	ticketID := ctx.Param("id")
	if ticketID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	// Simulate deleting the ticket in a database or service
	ctx.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully", "ticket_id": ticketID})
}

func GetTicketMessages(ctx *gin.Context) {
	ticketID := ctx.Param("id")
	if ticketID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	messages := []models.TicketMessage{}

	// Simulate fetching messages for the ticket from a database or service
	if err := ctx.ShouldBindQuery(&messages); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

func CreateTicketMessage(ctx *gin.Context) {}
