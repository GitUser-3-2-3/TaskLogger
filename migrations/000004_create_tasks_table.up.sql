CREATE TABLE tasks
(
    task_id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    status           task_status              DEFAULT 'Pending',
    priority         task_priority            DEFAULT 'medium',
    image_url        VARCHAR(500),
    duration_minutes INTEGER, -- Store as minutes for consistency
    deadline         TIMESTAMP WITH TIME ZONE,
    user_id          UUID         NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    category_id      UUID REFERENCES categories (ctg_id) ON DELETE CASCADE,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT positive_duration CHECK (duration_minutes > 0),
    CONSTRAINT future_deadline CHECK (deadline > created_at)
);

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_category_id_fkey;

ALTER TABLE tasks ALTER COLUMN category_id TYPE BIGINT USING NULL;

ALTER TABLE tasks ADD CONSTRAINT tasks_category_id_fkey
    FOREIGN KEY (category_id) REFERENCES categories (ctg_id) ON DELETE CASCADE;
