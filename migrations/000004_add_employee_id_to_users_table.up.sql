-- Add employee_id field to users table
-- Migration: 000004_add_employee_id_to_users_table.up.sql

ALTER TABLE users ADD COLUMN employee_id VARCHAR(50);

-- Create unique constraint on employee_id to ensure no duplicates
ALTER TABLE users ADD CONSTRAINT users_employee_id_unique UNIQUE (employee_id);

-- Create index for better query performance
CREATE INDEX idx_users_employee_id ON users(employee_id);

-- Add comment for documentation
COMMENT ON COLUMN users.employee_id IS 'Employee identifier (NIP) - unique identifier for employee';