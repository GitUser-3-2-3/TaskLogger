CREATE INDEX idx_tasks_user ON tasks (user_id);

CREATE INDEX idx_tasks_category ON tasks (category_id);

CREATE INDEX idx_categories_user ON categories (user_id);

CREATE INDEX idx_sessions_task ON sessions (task_id);
