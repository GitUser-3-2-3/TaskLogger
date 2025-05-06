CREATE TABLE IF NOT EXISTS tasks
(
    id             INTEGER AUTO_INCREMENT PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    description    TEXT,
    status         ENUM ('Not Started', 'In Progress', 'Paused', 'Completed') DEFAULT 'Not Started',
    priority       TINYINT                                                    DEFAULT 3,
    image          LONGBLOB,
    total_duration INTEGER                                                    DEFAULT 0,
    created_at     DATETIME     NOT NULL                                      DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME                                                   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deadline       DATETIME,
    user_id        INTEGER      NOT NULL,
    category_id    INTEGER,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL
);
