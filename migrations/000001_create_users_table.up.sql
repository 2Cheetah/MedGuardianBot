CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    first_name TEXT NULL,
    last_name TEXT NULL,
    username TEXT UNIQUE
);