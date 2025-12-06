CREATE TABLE activities (
    uuid VARCHAR(36) PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    parent_uuid VARCHAR(36),
    level INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255) NULL,
    FOREIGN KEY (parent_uuid) REFERENCES activities(uuid) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_activities_parent_uuid ON activities(parent_uuid);
CREATE INDEX idx_activities_level ON activities(level);
CREATE INDEX idx_activities_code ON activities(code);
CREATE INDEX idx_activities_deleted_at ON activities(deleted_at);
CREATE INDEX idx_activities_is_active ON activities(is_active);
