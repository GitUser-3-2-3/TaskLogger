CREATE TABLE sessions
(
    session_id UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    task_id    UUID                     NOT NULL REFERENCES tasks (task_id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at   TIMESTAMP WITH TIME ZONE,
    duration   INTEGER,
    notes      TEXT,
    created_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT valid_session_time CHECK (ended_at IS NULL OR ended_at > started_at)
);
