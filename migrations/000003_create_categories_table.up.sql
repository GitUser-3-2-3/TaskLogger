CREATE TABLE IF NOT EXISTS categories
(
    ctg_id     BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE NOT NULL,
    color      VARCHAR(7),
    user_id    UUID                NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_category_per_user UNIQUE (user_id, name)
);
