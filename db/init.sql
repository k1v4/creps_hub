CREATE TABLE "articles" (
                            "article_id" serial PRIMARY KEY,
                            "author_id" integer,
                            "publication_date" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            "name" varchar(100),
                            "text" text
);

CREATE TABLE "users" (
                         "id" serial PRIMARY KEY,
                         "email" varchar(50) UNIQUE,
                         "password" text,
                         "username" varchar(50) UNIQUE,
                         "name" text,
                         "surname" text,
                         "access_id" integer DEFAULT 1
);
CREATE INDEX idx_email_users ON users USING hash (email);
CREATE INDEX idx_username_users ON users USING hash (username);

CREATE TABLE "shoes" (
                         "id" serial PRIMARY KEY,
                         "user_id" integer,
                         "name" text,
                         "image_url" text
);

CREATE TABLE "releases" (
                            "id" serial PRIMARY KEY,
                            "date" timestamp,
                            "name" text
);

CREATE TABLE "access_levels" (
                                 "id" serial PRIMARY KEY,
                                 "level_name" varchar(20)
);

ALTER TABLE "users" ADD FOREIGN KEY ("access_id") REFERENCES "access_levels" ("id");

ALTER TABLE "articles" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");

ALTER TABLE "shoes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


INSERT INTO "access_levels"(level_name) VALUES ('user');


-- Создаем функцию, которая будет использоваться в триггере
-- CREATE OR REPLACE FUNCTION set_default_access_id()
--     RETURNS TRIGGER AS $$
-- BEGIN
--     -- Если access_id не указан, устанавливаем его равным 1
--     IF NEW.access_id IS NULL THEN
--         NEW.access_id := 1;
--     END IF;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- -- Создаем триггер, который будет вызывать функцию перед вставкой новой записи в таблицу users
-- CREATE TRIGGER trg_set_default_access_id
--     BEFORE INSERT ON users
--     FOR EACH ROW
-- EXECUTE FUNCTION set_default_access_id();