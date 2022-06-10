CREATE TABLE IF NOT EXISTS wn_message (
    user_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    group_id BIGINT REFERENCES wn_group(id) ON DELETE CASCADE,
    time_added TIMESTAMP NOT NULL,
    msg VARCHAR(512) NOT NULL,
    CHECK(msg != '')
);