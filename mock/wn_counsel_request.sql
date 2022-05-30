DROP TABLE wn_counsel_request;

CREATE TABLE wn_counsel_request (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT REFERENCES wn_user(id),
    counsel_title VARCHAR(256) NOT NULL,
    counsel_description VARCHAR(256) NOT NULL,
    request_status VARCHAR(8) NOT NULL,
    unique(user_id),
    check(counsel_title != ''),
    check(counsel_description != ''),
    check(request_status IN ('PENDING', 'REJECTED', 'APPROVED'))
)

GRANT ALL PRIVILEGES ON wn_counsel_request, wn_counsel_request_id_seq TO wellnus_admin;