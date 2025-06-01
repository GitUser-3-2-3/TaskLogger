DROP TABLE IF EXISTS tasks;

ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS positive_duration;

ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS future_deadline;
