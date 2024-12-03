CREATE TABLE IF NOT EXISTS user_data
(
    id         uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id    uuid        NOT NULL REFERENCES users (id),
    data_type  VARCHAR(70) NOT NULL,
    data       TEXT        NOT NULL,
    meta       TEXT,
    is_deleted bool DEFAULT false
)