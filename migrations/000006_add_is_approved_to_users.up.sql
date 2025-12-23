-- Add is_approved column to users table
-- Default to false for new registrations, requiring admin approval
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_approved BOOLEAN DEFAULT false;

-- Set existing users as approved (to not break existing users)
UPDATE users SET is_approved = true WHERE is_approved IS NULL;

-- Make the column NOT NULL after setting default values
ALTER TABLE users ALTER COLUMN is_approved SET NOT NULL;

-- Create index for filtering by approval status
CREATE INDEX IF NOT EXISTS idx_users_is_approved ON users(is_approved);
