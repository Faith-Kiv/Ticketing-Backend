CREATE TABLE tickets (
  id UUID PRIMARY KEY,
  customer_name VARCHAR(100) NOT NULL,
  customer_phone_number VARCHAR(30),
  source VARCHAR(30) NOT NULL,             -- 'email', 'chat', 'facebook', etc.
  subject VARCHAR(100) NOT NULL,
  content VARCHAR(500) NOT NULL,
  category VARCHAR(50) NOT NULL,  -- 'technical', 'billing', 'general', etc.
  priority VARCHAR(30),           -- 'low', 'medium', 'high'
  status VARCHAR(30) DEFAULT 'OPEN',
  agent_email VARCHAR(100),  -- Email of the agent assigned to the ticket
  agent_name VARCHAR(100),   -- Name of the agent assigned to the ticket
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  closed_at TIMESTAMP,
  resolved_at TIMESTAMP,
  sla_expires_at TIMESTAMP
);