CREATE TABLE IF NOT EXISTS wn_counsel_setting (
    user_id BIGINT PRIMARY KEY REFERENCES wn_user(id) ON DELETE CASCADE,
    available BOOLEAN NOT NULL
)