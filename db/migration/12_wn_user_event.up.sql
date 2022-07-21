CREATE TABLE IF NOT EXISTS wn_user_event (
    event_id BIGINT REFERENCES wn_event(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    unique(event_id, user_id)
)