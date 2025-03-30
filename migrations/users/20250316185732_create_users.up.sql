CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nickname TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password VARCHAR(100) NOT NULL
);
INSERT INTO users (id, nickname, email, hashed_password)
VALUES 
(1, 'sample', 'sample@sample.com', '$2a$12$B7HOO2geyYwV4T/uEt6lZ.eXlQl4rbxZKp.IIqef.1NhDlk9kUQE2');