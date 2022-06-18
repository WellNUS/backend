CREATE TABLE IF NOT EXISTS wn_session (
    session_key VARCHAR(128) NOT NULL PRIMARY KEY,
    user_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    unique(session_key),
    unique(user_id)
);