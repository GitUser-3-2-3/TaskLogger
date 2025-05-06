CREATE TABLE IF NOT EXISTS sessions
(
    id            INTEGER AUTO_INCREMENT PRIMARY KEY,
    task_id       INTEGER                NOT NULL,
    session_start DATETIME               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    session_end   DATETIME,
    duration      INTEGER,
    note          TEXT,
    session_type  ENUM ('work', 'break') NOT NULL DEFAULT 'work',
    FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
);
