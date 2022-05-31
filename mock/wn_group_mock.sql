DROP TABLE wn_user_group;
DROP TABLE wn_join_request;
DROP TABLE wn_group;

CREATE TABLE wn_group (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    group_name VARCHAR(64) NOT NULL,
    group_description VARCHAR(256) NOT NULL,
    category VARCHAR(7) NOT NULL,
    owner_id BIGINT REFERENCES wn_user(id) NOT NULL,
    check(group_name != ''),
    check(category IN ('COUNSEL', 'SUPPORT'))
);

GRANT ALL PRIVILEGES ON wn_group, wn_group_id_seq TO wellnus_admin;