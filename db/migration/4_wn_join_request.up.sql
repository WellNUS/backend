CREATE TABLE IF NOT EXISTS wn_join_request (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    group_id BIGINT REFERENCES wn_group(id) ON DELETE CASCADE,
    unique(user_id, group_id)
);