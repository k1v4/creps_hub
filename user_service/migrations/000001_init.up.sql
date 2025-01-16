CREATE TABLE IF NOT EXISTS users (
                                     id INT PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     surname TEXT NOT NULL,
                                     username TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_username ON users (username);

CREATE TABLE IF NOT EXISTS shoes (
                                     id SERIAL PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     image_url TEXT NOT NULL,
                                     user_id INT NOT NULL,
                                     FOREIGN KEY (user_id) REFERENCES users(id)
                                 );
