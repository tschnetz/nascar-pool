-- NASCAR Fantasy Pool Manager Database Schema
-- Run this file to create all necessary tables

-- Participants table (the 9 friends in the pool)
CREATE TABLE IF NOT EXISTS participants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Drivers table (NASCAR drivers with permanent car numbers)
CREATE TABLE IF NOT EXISTS drivers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    car_number INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Races table (simplified - results stored separately)
CREATE TABLE IF NOT EXISTS races (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    race_number INTEGER NOT NULL,
    date DATE,
    is_special_race BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'in_progress', 'completed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Race results table (tracks results by car number)
CREATE TABLE IF NOT EXISTS race_results (
    id SERIAL PRIMARY KEY,
    race_id INTEGER NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    car_number INTEGER NOT NULL,
    position INTEGER,  -- NULL for non-placing results
    is_first_place BOOLEAN DEFAULT FALSE,
    is_second_place BOOLEAN DEFAULT FALSE,
    is_last_place BOOLEAN DEFAULT FALSE,
    is_stage1_winner BOOLEAN DEFAULT FALSE,
    is_stage2_winner BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(race_id, car_number)
);

-- Teams table (tracks which drivers each participant has for each race)
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    race_id INTEGER NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    participant_id INTEGER NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    driver1_id INTEGER NOT NULL REFERENCES drivers(id),
    driver2_id INTEGER NOT NULL REFERENCES drivers(id),
    driver3_id INTEGER NOT NULL REFERENCES drivers(id),
    driver4_id INTEGER NOT NULL REFERENCES drivers(id),
    points_earned INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(race_id, participant_id)
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_teams_race_id ON teams(race_id);
CREATE INDEX IF NOT EXISTS idx_teams_participant_id ON teams(participant_id);
CREATE INDEX IF NOT EXISTS idx_races_status ON races(status);
CREATE INDEX IF NOT EXISTS idx_races_date ON races(date);
CREATE INDEX IF NOT EXISTS idx_race_results_race_id ON race_results(race_id);
CREATE INDEX IF NOT EXISTS idx_drivers_car_number ON drivers(car_number);

-- Seed initial participants
INSERT INTO participants (name) VALUES
    ('Foley'),
    ('Griffith'),
    ('Knight'),
    ('Patrick'),
    ('Rosing'),
    ('Schnetzer'),
    ('Simmons'),
    ('Wright'),
    ('Yancey')
ON CONFLICT (name) DO NOTHING;

-- Seed NASCAR drivers with car numbers
INSERT INTO drivers (name, car_number) VALUES
    ('Chastain', 1),
    ('Cindric', 2),
    ('Dillon A', 3),
    ('Gragson', 4),
    ('Larson', 5),
    ('Keselowski', 6),
    ('Haley', 7),
    ('Busch', 8),
    ('Elliott', 9),
    ('Dillon T', 10),
    ('Hamlin', 11),
    ('Blaney', 12),
    ('Allmendinger', 16),
    ('Buescher', 17),
    ('Briscoe', 14),
    ('Bell', 20),
    ('Berry', 21),     -- Verify car number for current season
    ('Logano', 22),
    ('Wallace', 23),
    ('Byron', 24),
    ('Gilliland', 38),
    ('Herbst', 35),
    ('Smith Z', 71),
    ('Custer', 41),
    ('Nemechek', 42),
    ('Jones', 43),
    ('Reddick', 45),
    ('Stenhouse', 47),
    ('Bowman', 48),
    ('Suarez', 99),
    ('McDowell', 34),
    ('Gibbs', 54),
    ('Preece', 33),
    ('Hocevar', 77),
    ('Van Gisbergen', 88),
    ('Ware', 51)
ON CONFLICT (name) DO NOTHING;
