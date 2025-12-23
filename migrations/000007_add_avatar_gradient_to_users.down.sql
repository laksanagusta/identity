-- Rollback avatar gradient columns from users table

DROP INDEX IF EXISTS idx_users_avatar_gradient;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_gradient_end;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_gradient_start;
