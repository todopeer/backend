ALTER TABLE events ADD COLUMN user_id INTEGER;
CREATE INDEX idx_events_user_id ON events(user_id);

-- currently we only have one user, fill in the empty field
UPDATE events SET user_id = 1 WHERE user_id IS NULL;

ALTER TABLE events ADD COLUMN description TEXT;