CREATE TABLE IF NOT EXISTS users
(
    id       uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    login    VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100)        NOT NULL
);