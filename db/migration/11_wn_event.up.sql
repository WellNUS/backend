CREATE TABLE IF NOT EXISTS wn_event (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    owner_id BIGINT REFERENCES wn_user(id) NOT NULL,
    event_name TEXT NOT NULL,
    event_description TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    access VARCHAR(7) NOT NULL,
    category VARCHAR(7) NOT NULL,
    check(event_name != ''),
    check(access IN ('PUBLIC', 'PRIVATE')),
    check(category IN ('COUNSEL', 'SUPPORT', 'CUSTOM'))
)