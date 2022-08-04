CREATE TABLE IF NOT EXISTS wn_direct (
    sender_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    recipient_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    time_added TIMESTAMPTZ NOT NULL,
    msg VARCHAR(512) NOT NULL,
    CHECK(msg != '')
);