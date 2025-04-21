CREATE TABLE "articles" (
                            "article_id" serial PRIMARY KEY,
                            "author_id" integer,
                            "publication_date" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            "name" varchar(100),
                            "text" text,
                            "image_url" text
);

CREATE TABLE "users" (
                         "id" serial PRIMARY KEY,
                         "email" varchar(50) UNIQUE,
                         "password" text,
                         "username" varchar(50) UNIQUE NOT NULL,
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
                            "date" timestamp UNIQUE,
                            "name" text,
                            "image_url" text
);

CREATE TABLE "access_levels" (
                                 "id" serial PRIMARY KEY,
                                 "level_name" varchar(20)
);

ALTER TABLE "users" ADD FOREIGN KEY ("access_id") REFERENCES "access_levels" ("id");
ALTER TABLE "articles" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");
ALTER TABLE "shoes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

INSERT INTO "users"(email, password, username) VALUES ('deleted@mail.ru', '1111111111', 'deleted');

INSERT INTO "access_levels"(level_name) VALUES ('user');
INSERT INTO "access_levels"(level_name) VALUES ('admin');
