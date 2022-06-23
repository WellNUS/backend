CREATE TABLE IF NOT EXISTS wn_match_setting (
    user_id BIGINT PRIMARY KEY REFERENCES wn_user(id) ON DELETE CASCADE,
    faculty_preference VARCHAR(4) NOT NULL,
    hobbies TEXT[],
    mbti VARCHAR(4) NOT NULL,
    CHECK(faculty_preference IN ('MIX', 'SAME', 'NONE')),
    CHECK(hobbies <@ ARRAY['GAMING', 'SINGING', 'DANCING', 'MUSIC', 'SPORTS', 'OUTDOOR', 'BOOK', 'ANIME', 'MOVIES', 'TV', 'ART', 'STUDY']),
    CHECK(array_length(hobbies, 1) <= 4),
    CHECK(mbti IN ('ISTJ','ISFJ','INFJ','INTJ','ISTP','ISFP','INFP','INTP','ESTP','ESFP','ENFP','ENTP','ESTJ','ESFJ','ENFJ','ENTJ'))
);