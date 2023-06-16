ALTER TABLE tasks ADD COLUMN deleted_at TIMESTAMP;
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);