CREATE TABLE IF NOT EXISTS wn_booking (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    recipient_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    provider_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    approve_by BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    nickname TEXT NOT NULL,
    details TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    unique(recipient_id, provider_id)
)