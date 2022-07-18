CREATE TABLE IF NOT EXISTS wn_group (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    group_name VARCHAR(64) NOT NULL,
    group_description VARCHAR(256) NOT NULL,
    category VARCHAR(7) NOT NULL,
    owner_id BIGINT REFERENCES wn_user(id) NOT NULL,
    check(group_name != ''),
    check(category IN ('COUNSEL', 'SUPPORT', 'CUSTOM'))
);