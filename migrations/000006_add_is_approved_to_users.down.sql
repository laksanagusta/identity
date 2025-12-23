-- Remove is_approved column from users table
DROP INDEX IF EXISTS idx_users_is_approved;
ALTER TABLE users DROP COLUMN IF EXISTS is_approved;
