CREATE TABLE IF NOT EXISTS wn_match_request (
    user_id BIGINT REFERENCES wn_match_setting(user_id) ON DELETE CASCADE,
    time_added TIMESTAMP NOT NULL,
    UNIQUE(user_id)
);