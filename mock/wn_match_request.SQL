DROP TABLE wn_match_request;

CREATE TABLE wn_match_request (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT REFERENCES wn_user(id),
    user_faculty VARCHAR(10) NOT NULL,
    cross_gender BOOLEAN NOT NULL,
    cross_faculty BOOLEAN NOT NULL,
    size SMALLINT NOT NULL, 
    mbti VARCHAR(4) NOT NULL,
    request_status VARCHAR(8) NOT NULL,
    unique(user_id),
    check(user_faculty IN ('CHS', 'BUSINESS', 'COMPUTING', 'DENTISTRY', 'CDE', 'LAW', 'MEDICINE', 'NURSING', 'PHARMACY', 'MUSIC')),
    check(size <= 8),
    check(mbti IN ('INFJ', 'INTP', 'ENTJ', 'ENTP', 'INFJ', 'INFP', 'ENFJ', 'ENFP', 'ISTJ', 'ISFJ', 'ESTJ', 'ESFJ', 'ISTP', 'ISFP', 'ESTP', 'ESFP'))
    check(request_status IN ('PENDING', 'REJECTED', 'APPROVED'))
)

GRANT ALL PRIVILEGES ON wn_match_request, wn_match_request_id_seq TO wellnus_admin;