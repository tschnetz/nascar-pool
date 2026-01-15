-- Seed 2026 NASCAR Cup Series Schedule
-- Delete any existing races first
DELETE FROM races;

-- Insert all 36 points races
-- Special races (2x points): Daytona 500, Coca-Cola 600, Southern 500, Championship
INSERT INTO races (race_number, name, date, is_special_race, status) VALUES
(1, 'DAYTONA 500', '2026-02-15', true, 'upcoming'),
(2, 'Ambetter Health 400', '2026-02-22', false, 'upcoming'),
(3, 'DuraMAX Grand Prix', '2026-03-01', false, 'upcoming'),
(4, 'Straight Talk Wireless 500', '2026-03-08', false, 'upcoming'),
(5, 'Pennzoil 400', '2026-03-15', false, 'upcoming'),
(6, 'Goodyear 400', '2026-03-22', false, 'upcoming'),
(7, 'Cook Out 400', '2026-03-29', false, 'upcoming'),
(8, 'Food City 500', '2026-04-12', false, 'upcoming'),
(9, 'AdventHealth 400', '2026-04-19', false, 'upcoming'),
(10, 'Jack Link''s 500', '2026-04-26', false, 'upcoming'),
(11, 'Würth 400', '2026-05-03', false, 'upcoming'),
(12, 'Go Bowling at the Glen', '2026-05-10', false, 'upcoming'),
(13, 'Coca-Cola 600', '2026-05-24', true, 'upcoming'),
(14, 'Cracker Barrel 400', '2026-05-31', false, 'upcoming'),
(15, 'FireKeepers Casino 400', '2026-06-07', false, 'upcoming'),
(16, 'Great American Getaway 400', '2026-06-14', false, 'upcoming'),
(17, 'Anduril 250 Race the Base', '2026-06-21', false, 'upcoming'),
(18, 'Toyota / Save Mart 350', '2026-06-28', false, 'upcoming'),
(19, 'TBA 400', '2026-07-05', false, 'upcoming'),
(20, 'Quaker State 400', '2026-07-12', false, 'upcoming'),
(21, 'Window World 450', '2026-07-19', false, 'upcoming'),
(22, 'Brickyard 400', '2026-07-26', false, 'upcoming'),
(23, 'Iowa Corn 350', '2026-08-09', false, 'upcoming'),
(24, 'Cook Out 400', '2026-08-15', false, 'upcoming'),
(25, 'USA Today 301', '2026-08-23', false, 'upcoming'),
(26, 'Coke Zero Sugar 400', '2026-08-29', false, 'upcoming'),
(27, 'Cook Out Southern 500', '2026-09-06', true, 'upcoming'),
(28, 'Enjoy Illinois 300', '2026-09-13', false, 'upcoming'),
(29, 'Bass Pro Shops Night Race', '2026-09-19', false, 'upcoming'),
(30, 'Hollywood Casino 400', '2026-09-27', false, 'upcoming'),
(31, 'South Point 400', '2026-10-04', false, 'upcoming'),
(32, 'Bank of America ROVAL 400', '2026-10-11', false, 'upcoming'),
(33, 'TBA', '2026-10-18', false, 'upcoming'),
(34, 'YellaWood 500', '2026-10-25', false, 'upcoming'),
(35, 'Xfinity 500', '2026-11-01', false, 'upcoming'),
(36, 'NASCAR Cup Series Championship', '2026-11-08', true, 'upcoming');
