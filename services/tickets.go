package services

import (
	"fmt"
	"time"

	"github.com/Faith-Kiv/Ticketing-Backend/utils"
	"github.com/nyaruka/phonenumbers"
	"github.com/sirupsen/logrus"
)

func GetTickets() string {
	// This function is a placeholder for the ticket service logic.
	// It can be expanded to include actual service logic as needed.
	return "Ticket Service is running"
}

func CreateTicket() string {
	// This function is a placeholder for the ticket creation logic.
	// It can be expanded to include actual service logic as needed.
	return "Ticket created successfully"
}

func GetTicket() string {
	// This function is a placeholder for retrieving a specific ticket.
	// It can be expanded to include actual service logic as needed.
	return "Ticket details retrieved successfully"
}

func UpdateTicket() string {
	// This function is a placeholder for updating a specific ticket.
	// It can be expanded to include actual service logic as needed.
	return "Ticket updated successfully"
}

func DeleteTicket() string {
	// This function is a placeholder for deleting a specific ticket.
	// It can be expanded to include actual service logic as needed.
	return "Ticket deleted successfully"
}

func GetTicketMessages() string {
	// This function is a placeholder for retrieving messages associated with a ticket.
	// It can be expanded to include actual service logic as needed.
	return "Ticket messages retrieved successfully"
}

func CreateTicketMessage() string {
	// This function is a placeholder for creating a message associated with a ticket.
	// It can be expanded to include actual service logic as needed.
	return "Ticket message created successfully"
}

func ValidatePhoneNumber(phoneNumber string) bool {
	// This function is a placeholder for phone number validation logic.
	// It can be expanded to include actual validation logic as needed.
	parsed, err := phonenumbers.Parse(phoneNumber, utils.DEFAULT_COUNTRY_CODE)
	if err != nil || !phonenumbers.IsValidNumber(parsed) {
		logrus.Error(err)
		err = fmt.Errorf("invalid phone number %s", phoneNumber)
		return false
	}
	return true
}

func CalculateSLA(priority string) time.Time {
	now := time.Now()
	switch priority {
	case "high":
		return now.Add(2 * time.Hour)
	case "medium":
		return now.Add(4 * time.Hour)
	default:
		return now.Add(12 * time.Hour)
	}
}
