
CREATE UNIQUE INDEX idx_phone_key on users (phone) where deleted_at IS NULL;
CREATE UNIQUE INDEX idx_email_key on users (email) where deleted_at IS NULL;
