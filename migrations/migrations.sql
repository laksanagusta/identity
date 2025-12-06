CREATE TABLE activities (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMP,
    deleted_by TEXT,

    code TEXT,
    name TEXT,
    type TEXT,
    parent_uuid UUID,
    level INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_activity_parent FOREIGN KEY (parent_uuid) REFERENCES activity(uuid)
);

CREATE TABLE budget_allocations (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMP,
    deleted_by TEXT,

    fiscal_year TEXT,
    organization_uuid TEXT,
    program_uuid UUID,
    activity_uuid UUID,
    funds_type TEXT,
    amount TEXT,

    CONSTRAINT fk_budget_allocation_program FOREIGN KEY (program_uuid) REFERENCES activity(uuid),
    CONSTRAINT fk_budget_allocation_activity FOREIGN KEY (activity_uuid) REFERENCES activity(uuid)
);

