CREATE TABLE IF NOT EXISTS schedule_slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    start_time VARCHAR(5) NOT NULL,
    end_time VARCHAR(5) NOT NULL,
    title VARCHAR(200) DEFAULT '',
    recurring BOOLEAN DEFAULT true,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_schedule_creator_id ON schedule_slots(creator_id);
CREATE INDEX IF NOT EXISTS idx_schedule_day_time ON schedule_slots(day_of_week, start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_schedule_active ON schedule_slots(creator_id, active);
CREATE INDEX IF NOT EXISTS idx_schedule_live_now ON schedule_slots(day_of_week, start_time, end_time, active);

CREATE UNIQUE INDEX IF NOT EXISTS idx_schedule_no_overlap 
ON schedule_slots(creator_id, day_of_week, start_time) 
WHERE active = true;
