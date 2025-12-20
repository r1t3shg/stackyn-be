-- Add slug and user_id columns to apps table
ALTER TABLE apps 
ADD COLUMN IF NOT EXISTS slug VARCHAR(255),
ADD COLUMN IF NOT EXISTS user_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS status VARCHAR(255),
ADD COLUMN IF NOT EXISTS url VARCHAR(255),
ADD COLUMN IF NOT EXISTS branch VARCHAR(255);

-- Create index on user_id for faster queries
CREATE INDEX IF NOT EXISTS idx_apps_user_id ON apps(user_id);


