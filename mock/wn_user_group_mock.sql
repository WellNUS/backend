DROP TABLE wn_user_group;

CREATE TABLE wn_user_group (
    user_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    group_id BIGINT REFERENCES wn_group(id) ON DELETE CASCADE,
    unique(user_id, group_id)
);

GRANT ALL PRIVILEGES ON wn_user_group TO wellnus_admin;