CREATE TABLE IF NOT EXISTS users (
                                     id INT PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     surname TEXT NOT NULL,
                                     username TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_username ON users (username);
