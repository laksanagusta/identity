CREATE TABLE events (
    uuid VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    workflow_uuid VARCHAR(36) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    start_date TIMESTAMP NULL,
    end_date TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255) NULL,
    FOREIGN KEY (workflow_uuid) REFERENCES workflows(uuid) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_events_workflow_uuid ON events(workflow_uuid);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_start_date ON events(start_date);
CREATE INDEX idx_events_end_date ON events(end_date);
CREATE INDEX idx_events_deleted_at ON events(deleted_at);
