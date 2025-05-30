CREATE INDEX idx_tasks_user_id ON tasks (user_id);
CREATE INDEX idx_tasks_status ON tasks (status);
CREATE INDEX idx_tasks_priority ON tasks (priority);
CREATE INDEX idx_tasks_deadline ON tasks (deadline) WHERE deadline IS NOT NULL;
CREATE INDEX idx_tasks_category_id ON tasks (category_id);

CREATE INDEX idx_categories_user_id ON categories (user_id);

CREATE INDEX idx_sessions_task_id ON sessions (task_id);
CREATE INDEX idx_sessions_started_at ON sessions (started_at);
CREATE INDEX idx_sessions_active ON sessions (task_id) WHERE ended_at IS NULL;
