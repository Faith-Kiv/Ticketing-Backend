package models

import "bytes"

type Tickets struct {
	ID                  string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CustomerName        string `json:"customer_name" gorm:"not null;size:100"`
	CustomerPhoneNumber string `json:"customer_phone_number" gorm:"size:30"`
	Source              string `json:"source" gorm:"not null;size:30"` // 'email', 'chat', 'facebook', etc.
	Subject             string `json:"subject" gorm:"not null;size:100"`
	Content             string `json:"content" gorm:"not null;size:500"`
	Category            string `json:"category" gorm:"not null;size:50"`     // 'technical', 'billing', 'general', etc.
	Priority            string `json:"priority" gorm:"size:30"`              // 'low', 'medium', 'high'
	Status              string `json:"status" gorm:"default:'OPEN';size:30"` // 'OPEN', 'CLOSED', 'PENDING', etc.
	AgentEmail          string `json:"agent_email" gorm:"size:100"`          // Email of the agent assigned to the ticket
	AgentName           string `json:"agent_name" gorm:"size:100"`           // Name of the agent assigned to the ticket
	CreatedAt           string `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt           string `json:"updated_at" gorm:"default:current_timestamp"`
	ClosedAt            string `json:"closed_at" gorm:"default:null"`
	ResolvedAt          string `json:"resolved_at" gorm:"default:null"`
	SlaExpiresAt        string `json:"sla_expires_at" gorm:"default:null"`
}

type TicketMessage struct {
	ID          string      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   string      `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   string      `json:"updated_at" gorm:"default:current_timestamp"`
	Author      string      `json:"author" gorm:"not null;size:100"` // 'customer' or 'agent'
	Content     string      `json:"content" gorm:"not null;size:500"`
	Attachments []InMemFile `json:"attachments" gorm:"-"`
	// Attachments are not stored in the database, but can be included in the response
	// They are handled separately in the application logic.
}

type InMemFile struct {
	FileName string
	Buffer   bytes.Buffer
}
