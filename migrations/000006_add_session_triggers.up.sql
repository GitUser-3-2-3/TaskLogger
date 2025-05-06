DELIMITER //
CREATE TRIGGER update_task_duration_after_session
    AFTER INSERT
    ON sessions
    FOR EACH ROW
BEGIN
    IF NEW.duration IS NOT NULL THEN
        UPDATE tasks
        SET total_duration = total_duration + NEW.duration,
            updated_at     = CURRENT_TIMESTAMP
        WHERE id = NEW.task_id;
    END IF;
END;
//
DELIMITER ;

DELIMITER //
CREATE TRIGGER update_task_after_session
    AFTER INSERT
    ON sessions
    FOR EACH ROW
BEGIN
    IF NEW.session_type = 'work' AND NEW.duration IS NULL THEN
        UPDATE tasks SET status = 'In Progress', updated_at = CURRENT_TIMESTAMP WHERE id = NEW.task_id;
    ELSEIF NEW.session_type = 'work' AND NEW.duration IS NOT NULL THEN
        UPDATE tasks SET status = 'Completed', updated_at = CURRENT_TIMESTAMP WHERE id = NEW.task_id;
    ELSEIF NEW.session_type = 'break' THEN
        UPDATE tasks SET status = 'Paused', updated_at = CURRENT_TIMESTAMP WHERE id = NEW.task_id;
    END IF;
END;
//
DELIMITER ;
