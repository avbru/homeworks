CREATE TABLE events (
event_id serial PRIMARY KEY,
title TEXT NOT NULL,
description TEXT,
start_time TIMESTAMP NOT NULL,
end_time TIMESTAMP NOT NULL,
notify_time TIMESTAMP NOT NULL,
user_id TEXT NOT NULL,
CHECK (end_time >= start_time),
EXCLUDE USING gist (tsrange(start_time, end_time) WITH &&)
);

-- ID string
-- Title string
-- StartTime time.Time
-- EndTime time.Time
-- Description string
-- UserID string
-- NotifyBefore time.Duration