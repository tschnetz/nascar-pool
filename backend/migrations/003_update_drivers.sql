-- Add team and manufacturer columns to drivers table
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS team_name VARCHAR(100);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS manufacturer VARCHAR(20);

-- Clear existing drivers and insert the 36 chartered drivers for 2026
DELETE FROM drivers;

INSERT INTO drivers (name, car_number, team_name, manufacturer) VALUES
('Ross Chastain', 1, 'Trackhouse Racing', 'Chevrolet'),
('Austin Cindric', 2, 'Team Penske', 'Ford'),
('Austin Dillon', 3, 'Richard Childress Racing', 'Chevrolet'),
('Noah Gragson', 4, 'Front Row Motorsports', 'Ford'),
('Kyle Larson', 5, 'Hendrick Motorsports', 'Chevrolet'),
('Brad Keselowski', 6, 'Roush Fenway Keselowski Racing', 'Ford'),
('Daniel Suarez', 7, 'Spire Motorsports', 'Chevrolet'),
('Kyle Busch', 8, 'Richard Childress Racing', 'Chevrolet'),
('Chase Elliott', 9, 'Hendrick Motorsports', 'Chevrolet'),
('Ty Dillon', 10, 'Kaulig Racing', 'Chevrolet'),
('Denny Hamlin', 11, 'Joe Gibbs Racing', 'Toyota'),
('Ryan Blaney', 12, 'Team Penske', 'Ford'),
('AJ Allmendinger', 16, 'Kaulig Racing', 'Chevrolet'),
('Chris Buescher', 17, 'Roush Fenway Keselowski Racing', 'Ford'),
('Chase Briscoe', 19, 'Joe Gibbs Racing', 'Toyota'),
('Christopher Bell', 20, 'Joe Gibbs Racing', 'Toyota'),
('Josh Berry', 21, 'Wood Brothers Racing', 'Ford'),
('Joey Logano', 22, 'Team Penske', 'Ford'),
('Bubba Wallace', 23, '23XI Racing', 'Toyota'),
('William Byron', 24, 'Hendrick Motorsports', 'Chevrolet'),
('Todd Gilliland', 34, 'Front Row Motorsports', 'Ford'),
('Riley Herbst', 35, '23XI Racing', 'Toyota'),
('Zane Smith', 38, 'Front Row Motorsports', 'Ford'),
('Cole Custer', 41, 'Haas Factory Team', 'Chevrolet'),
('John Hunter Nemechek', 42, 'Legacy Motor Club', 'Toyota'),
('Erik Jones', 43, 'Legacy Motor Club', 'Toyota'),
('Tyler Reddick', 45, '23XI Racing', 'Toyota'),
('Ricky Stenhouse Jr.', 47, 'HYAK Motorsports', 'Chevrolet'),
('Alex Bowman', 48, 'Hendrick Motorsports', 'Chevrolet'),
('Cody Ware', 51, 'Rick Ware Racing', 'Chevrolet'),
('Ty Gibbs', 54, 'Joe Gibbs Racing', 'Toyota'),
('Ryan Preece', 60, 'Roush Fenway Keselowski', 'Ford'),
('Michael McDowell', 71, 'Spire Motorsports', 'Chevrolet'),
('Carson Hocevar', 77, 'Spire Motorsports', 'Chevrolet'),
('Connor Zilisch', 88, 'Trackhouse Racing', 'Chevrolet'),
('Shane van Gisbergen', 97, 'Trackhouse Racing', 'Chevrolet');
