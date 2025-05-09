INSERT INTO users
    (username, email, firstname, lastname, password_hash)
VALUES ('parth', 'parth@gmail.com', 'Parth', 'Srivastava', 'labalabadubdub');

INSERT INTO users
    (username, email, firstname, lastname, password_hash)
VALUES ('juhi', 'juhi@gmail.com', 'Juhi', 'Srivastava', 'labalabadubdub');

SELECT * FROM users;
SELECT * FROM tasks;
SELECT id, name, color FROM categories WHERE user_id = 1;
