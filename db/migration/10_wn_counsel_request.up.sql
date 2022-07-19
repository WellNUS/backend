CREATE TABLE IF NOT EXISTS wn_counsel_request (
    user_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE,
    details TEXT,
    topics TEXT[] NOT NULL,
    lastUpdated TIMESTAMP NOT NULL,
    unique(user_id),
    check(COALESCE(array_length(topics, 1), 0) > 0),
    check(topics <@ ARRAY['Anxiety', 'OffMyChest', 'SelfHarm', 'Depression', 'SelfEsteem', 'Stress', 'Casual', 'Therapy', 'BadHabits', 'Rehabilitation'])
)