package main

import (
	"github.com/Faith-Kiv/Ticketing-Backend/controllers"
	"github.com/gin-gonic/gin"
)

var Routes = map[string]map[string]gin.HandlerFunc{
	"/api/tickets": {
		"GET":  controllers.GetTickets,
		"POST": controllers.CreateTicket,
	},

	"/api/ticket/:id": {
		"GET": controllers.GetTicket,
		"PUT": controllers.UpdateTicket,
	},

	"/api/ticket/:id/messages": {
		"GET":  controllers.GetTicketMessages,
		"POST": controllers.CreateTicketMessage,
	},
}
