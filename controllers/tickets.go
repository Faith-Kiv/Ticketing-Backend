package controllers

import (
	"net/http"
	"slices"

	"github.com/Faith-Kiv/Ticketing-Backend/database"
	"github.com/Faith-Kiv/Ticketing-Backend/models"
	"github.com/Faith-Kiv/Ticketing-Backend/services"
	gormabs "github.com/ecoa-dev-team/burngormabs"
	"github.com/gin-gonic/gin"
)

func GetTickets(ctx *gin.Context) {
	tickets := []models.Tickets{}
	params := ctx.Request.URL.Query()
	query := database.DbR

	// Simulate fetching tickets from a database or service
	if err := ctx.ShouldBindQuery(&tickets); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	err := gormabs.SearchMulti(params, query, models.Tickets{}, &tickets)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	rows, err := gormabs.Count(params, query, models.Tickets{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to count tickets"})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"count": rows, "ticket": tickets})
}

func CreateTicket(ctx *gin.Context) {
	roles := ctx.GetStringSlice("roles")
	email := ctx.GetString("email")
	name := ctx.GetString("username")
	ticket := models.Tickets{}

	if err := ctx.ShouldBindJSON(&ticket); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if !slices.Contains(roles, "customer_support") {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: You do not have permission to create a ticket."})
		return
	}

	if !services.ValidatePhoneNumber(ticket.CustomerPhoneNumber) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}

	// === AI CLASSIFICATION ===
	classification := services.ClassifyTicketAI(ticket.Content)
	ticket.Category = classification.Category
	ticket.Priority = classification.Priority
	ticket.Summary = classification.Summary

	// === AGENT ROUTING ===
	ticket.AgentEmail = email
	ticket.AgentName = name

	// === SLA DEADLINE ===
	ticket.SlaExpiresAt = services.CalculateSLA(ticket.Priority)

	if err := database.DbR.Table(models.Tickets{}.TableName()).Create(&ticket).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ticket"})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusCreated, ticket)

}

func GetTicket(ctx *gin.Context) {
	ticketID := ctx.Param("id")
	ticket := models.Tickets{}
	query := database.DbR

	if ticketID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	param := map[string][]string{"eq__id": {ticketID}}

	err := gormabs.SearchOne(param, query, &ticket)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, ticket)
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

	if err := database.DbR.Table(models.Tickets{}.TableName()).Where("id = ?", ticketID).Delete(&models.Tickets{}).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
		return
	}

	// Simulate deleting the ticket in a database or service
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully", "ticket_id": ticketID})
}

func GetTicketMessages(ctx *gin.Context) {
	ticketID := ctx.Param("id")
	messages := []models.TicketMessage{}
	query := database.DbR

	if ticketID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	param := map[string][]string{"eq__id": {ticketID}}
	err := gormabs.SearchMulti(param, query, models.Tickets{}, &messages)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	rows, err := gormabs.Count(param, query, models.TicketMessage{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to count tickets"})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"count": rows, "messages": messages})
}

func CreateTicketMessage(ctx *gin.Context) {
	ticketID := ctx.Param("id")
	email := ctx.GetString("email")
	message := models.TicketMessage{}

	if ticketID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	message.Author = email // Assuming the author is the agent's email

	if err := database.DbR.Table(models.TicketMessage{}.TableName()).Create(&message).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusCreated, message)
}
