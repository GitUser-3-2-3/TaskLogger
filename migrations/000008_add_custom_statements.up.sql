INSERT INTO users (username, email, firstname, lastname, password_hash)
VALUES ('parth', 'parth@gmail.com', 'Parth', 'Srivastava', 'labalabadubdub');

INSERT INTO users (username, email, firstname, lastname, password_hash)
VALUES ('juhi', 'juhi@gmail.com', 'Juhi', 'Srivastava', 'labalabadubdub');

SELECT *
FROM users;

SELECT *
FROM tasks;

SELECT *
FROM categories;

SELECT *
FROM sessions;

DELETE
FROM tasks
WHERE id = 2;

SELECT id, name, color
FROM categories
WHERE user_id = 1;

DESCRIBE tasks;
SELECT id, name, description, status, priority, image, total_duration, created_at, updated_at, deadline, category_id
FROM tasks
WHERE id = ?;
