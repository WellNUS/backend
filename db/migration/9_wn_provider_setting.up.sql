CREATE TABLE IF NOT EXISTS wn_provider_setting (
    user_id BIGINT PRIMARY KEY REFERENCES wn_user(id) ON DELETE CASCADE,
    available BOOLEAN NOT NULL,
    specialites TEXT[],
    unique(user_id),
    check(array_length(specialites, 1) > 0),
    check(specialites <@ ARRAY['Anxiety', 'OffMyChest', 'SelfHarm', 'Depression', 'SelfEsteem', 'Stress', 'Casual', 'Therapy', 'BadHabits', 'Rehabilitation'])
)