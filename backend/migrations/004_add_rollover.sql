-- Add rollover tracking for each scoring position
-- These track points rolling INTO this race from previous unassigned scorers

ALTER TABLE races ADD COLUMN IF NOT EXISTS rollover_first INTEGER DEFAULT 0;
ALTER TABLE races ADD COLUMN IF NOT EXISTS rollover_second INTEGER DEFAULT 0;
ALTER TABLE races ADD COLUMN IF NOT EXISTS rollover_last INTEGER DEFAULT 0;
ALTER TABLE races ADD COLUMN IF NOT EXISTS rollover_stage1 INTEGER DEFAULT 0;
ALTER TABLE races ADD COLUMN IF NOT EXISTS rollover_stage2 INTEGER DEFAULT 0;

-- Track which car numbers are chartered (the 36 that get assigned to teams)
-- Non-chartered cars that score will trigger rollover
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS is_chartered BOOLEAN DEFAULT true;
