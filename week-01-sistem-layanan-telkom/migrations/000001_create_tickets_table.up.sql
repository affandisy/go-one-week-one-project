CREATE TABLE tickets (
    ticket_id VARCHAR(50) PRIMARY KEY,
    customer_id VARCHAR(50) NOT NULL,
    issue TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'open',
    assigned_agent VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);