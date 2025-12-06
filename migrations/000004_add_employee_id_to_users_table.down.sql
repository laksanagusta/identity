-- Remove employee_id field from users table
-- Migration: 000004_add_employee_id_to_users_table.down.sql

-- Drop the index first
DROP INDEX IF EXISTS idx_users_employee_id;

-- Remove the unique constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_employee_id_unique;

-- Remove the column
ALTER TABLE users DROP COLUMN IF EXISTS employee_id;