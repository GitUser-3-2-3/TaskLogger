ALTER TABLE tasks
    ADD CONSTRAINT priority_level_check CHECK ( priority BETWEEN 1 AND 5);
