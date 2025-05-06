CREATE TABLE IF NOT EXISTS categories
(
    id      INTEGER AUTO_INCREMENT PRIMARY KEY,
    name    VARCHAR(50) NOT NULL,
    color   VARCHAR(7) DEFAULT '#FFFFFF',
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    UNIQUE (user_id, name)
)
