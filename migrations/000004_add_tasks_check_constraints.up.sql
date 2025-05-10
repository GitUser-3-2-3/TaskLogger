ALTER TABLE tasks
    ADD CONSTRAINT priority_level_check CHECK ( priority BETWEEN 1 AND 5);

ALTER TABLE tasks
    ADD CONSTRAINT tasks_ibfk_2
        FOREIGN KEY (category_id) REFERENCES categories (id)
            ON DELETE CASCADE;
