CREATE TABLE tickets (
  id UUID PRIMARY KEY,
  customer_name TEXT,
  source TEXT,             -- 'email', 'chat', 'facebook', etc.
  subject TEXT,
  content TEXT,
  category TEXT,
  priority TEXT,           -- 'low', 'medium', 'high'
  summary TEXT,
  status TEXT DEFAULT 'open',
  agent_id UUID REFERENCES agents(id),
  created_at TIMESTAMP DEFAULT NOW(),
  sla_expires_at TIMESTAMP
);