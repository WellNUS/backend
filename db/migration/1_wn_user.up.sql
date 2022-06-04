CREATE TABLE IF NOT EXISTS wn_user (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    first_name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL,
    gender VARCHAR(1) NOT NULL,
    faculty VARCHAR(10) NOT NULL,
    email VARCHAR(128) NOT NULL,
    user_role VARCHAR(10) NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    unique(email),
    check(first_name != ''),
    check(last_name != ''),
    check(password_hash != ''),
    check(gender IN ('M', 'F')),
    check(faculty IN ('CHS', 'BUSINESS', 'COMPUTING', 'DENTISTRY', 'CDE', 'LAW', 'MEDICINE', 'NURSING', 'PHARMACY', 'MUSIC')),
    check(email LIKE '%@u.nus.edu'),
    check(user_role IN ('MEMBER', 'VOLUNTEER', 'COUNSELLOR'))
);