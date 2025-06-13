package services

import (
	"strings"

	"github.com/Faith-Kiv/Ticketing-Backend/models"
)

func ClassifyTicketAI(content string) models.ClassificationResult {
	if strings.Contains(content, "bill") {
		return models.ClassificationResult{"billing", "medium", "Billing issue"}
	}
	return models.ClassificationResult{"general", "low", "General request"}
}
