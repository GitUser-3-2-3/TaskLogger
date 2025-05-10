ALTER TABLE tasks
    DROP CONSTRAINT priority_level_check;

ALTER TABLE tasks
    MODIFY COLUMN category_id INTEGER,
    DROP FOREIGN KEY tasks_ibfk_2;
