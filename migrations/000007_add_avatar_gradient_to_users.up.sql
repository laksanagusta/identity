-- Add avatar gradient columns to users table
-- These columns store the gradient colors for default profile photo

ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_gradient_start VARCHAR(7);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_gradient_end VARCHAR(7);

-- Create index for potential filtering if needed
CREATE INDEX IF NOT EXISTS idx_users_avatar_gradient ON users(avatar_gradient_start, avatar_gradient_end);
