-- Seed data for development

-- ============================================
-- Sports
-- ============================================
INSERT INTO sports (name, slug, icon) VALUES
    ('Football', 'football', '⚽'),
    ('Basketball', 'basketball', '🏀'),
    ('Tennis', 'tennis', '🎾');

-- ============================================
-- Admin user (password: admin123)
-- Hash generated with bcrypt cost 12
-- ============================================
INSERT INTO users (email, password_hash, balance, role) VALUES
    ('admin@betkz.com', '$2a$12$LJ3m4ys3LnFDNgdRAd4qJeXMxUbGCvVClHcuPA7AJJXSVf6OO2dCO', 10000.00, 'admin');

-- ============================================
-- Test user (password: user1234)
-- ============================================
INSERT INTO users (email, password_hash, balance, role) VALUES
    ('user@betkz.com', '$2a$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 500.00, 'user');

-- ============================================
-- Events (upcoming matches)
-- ============================================
INSERT INTO events (sport_id, home_team, away_team, start_time, status) VALUES
    (1, 'Real Madrid', 'Barcelona', NOW() + INTERVAL '2 days', 'upcoming'),
    (1, 'Manchester City', 'Liverpool', NOW() + INTERVAL '3 days', 'upcoming'),
    (1, 'Bayern Munich', 'Borussia Dortmund', NOW() + INTERVAL '1 day', 'upcoming'),
    (1, 'PSG', 'Marseille', NOW() + INTERVAL '4 days', 'upcoming'),
    (2, 'LA Lakers', 'Golden State Warriors', NOW() + INTERVAL '1 day', 'upcoming'),
    (2, 'Boston Celtics', 'Miami Heat', NOW() + INTERVAL '2 days', 'upcoming'),
    (2, 'Chicago Bulls', 'Brooklyn Nets', NOW() + INTERVAL '5 days', 'upcoming'),
    (3, 'Novak Djokovic', 'Carlos Alcaraz', NOW() + INTERVAL '3 days', 'upcoming'),
    (3, 'Rafael Nadal', 'Roger Federer', NOW() + INTERVAL '2 days', 'upcoming'),
    (1, 'Juventus', 'AC Milan', NOW() + INTERVAL '6 days', 'upcoming');

-- ============================================
-- Markets & Odds for first 5 events
-- ============================================

-- Event 1: Real Madrid vs Barcelona
-- 1X2 Market
INSERT INTO markets (event_id, market_type, name, status) 
SELECT id, '1x2', 'Match Winner', 'open' FROM events WHERE home_team = 'Real Madrid' AND away_team = 'Barcelona';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'home', 2.10, 2.10 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Real Madrid' AND e.away_team = 'Barcelona' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'draw', 3.40, 3.40 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Real Madrid' AND e.away_team = 'Barcelona' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'away', 3.20, 3.20 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Real Madrid' AND e.away_team = 'Barcelona' AND m.market_type = '1x2';

-- Over/Under Market
INSERT INTO markets (event_id, market_type, name, line, status)
SELECT id, 'over_under', 'Over/Under 2.5 Goals', 2.5, 'open' FROM events WHERE home_team = 'Real Madrid' AND away_team = 'Barcelona';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'over', 1.85, 1.85 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Real Madrid' AND e.away_team = 'Barcelona' AND m.market_type = 'over_under'
UNION ALL
SELECT m.id, 'under', 1.95, 1.95 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Real Madrid' AND e.away_team = 'Barcelona' AND m.market_type = 'over_under';

-- Both Teams Score
INSERT INTO markets (event_id, market_type, name, status)
SELECT id, 'both_teams_score', 'Both Teams to Score', 'open' FROM events WHERE home_team = 'Real Madrid' AND away_team = 'Barcelona';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'yes', 1.70, 1.70 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Real Madrid' AND e.away_team = 'Barcelona' AND m.market_type = 'both_teams_score'
UNION ALL
SELECT m.id, 'no', 2.10, 2.10 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Real Madrid' AND e.away_team = 'Barcelona' AND m.market_type = 'both_teams_score';

-- Event 2: Manchester City vs Liverpool (1X2 + Over/Under)
INSERT INTO markets (event_id, market_type, name, status)
SELECT id, '1x2', 'Match Winner', 'open' FROM events WHERE home_team = 'Manchester City' AND away_team = 'Liverpool';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'home', 1.90, 1.90 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Manchester City' AND e.away_team = 'Liverpool' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'draw', 3.60, 3.60 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Manchester City' AND e.away_team = 'Liverpool' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'away', 3.50, 3.50 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Manchester City' AND e.away_team = 'Liverpool' AND m.market_type = '1x2';

INSERT INTO markets (event_id, market_type, name, line, status)
SELECT id, 'over_under', 'Over/Under 2.5 Goals', 2.5, 'open' FROM events WHERE home_team = 'Manchester City' AND away_team = 'Liverpool';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'over', 1.75, 1.75 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Manchester City' AND e.away_team = 'Liverpool' AND m.market_type = 'over_under'
UNION ALL
SELECT m.id, 'under', 2.05, 2.05 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Manchester City' AND e.away_team = 'Liverpool' AND m.market_type = 'over_under';

-- Event 3: Bayern Munich vs Dortmund (1X2)
INSERT INTO markets (event_id, market_type, name, status)
SELECT id, '1x2', 'Match Winner', 'open' FROM events WHERE home_team = 'Bayern Munich';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'home', 1.65, 1.65 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Bayern Munich' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'draw', 4.00, 4.00 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Bayern Munich' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'away', 4.50, 4.50 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Bayern Munich' AND m.market_type = '1x2';

-- Event 5: Lakers vs Warriors (1X2 basketball)
INSERT INTO markets (event_id, market_type, name, status)
SELECT id, '1x2', 'Match Winner', 'open' FROM events WHERE home_team = 'LA Lakers';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'home', 1.80, 1.80 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'LA Lakers' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'away', 2.00, 2.00 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'LA Lakers' AND m.market_type = '1x2';

-- Event 8: Djokovic vs Alcaraz (1X2 tennis)
INSERT INTO markets (event_id, market_type, name, status)
SELECT id, '1x2', 'Match Winner', 'open' FROM events WHERE home_team = 'Novak Djokovic';

INSERT INTO odds (market_id, outcome, initial_odds, current_odds)
SELECT m.id, 'home', 1.55, 1.55 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Novak Djokovic' AND m.market_type = '1x2'
UNION ALL
SELECT m.id, 'away', 2.40, 2.40 FROM markets m JOIN events e ON m.event_id = e.id WHERE e.home_team = 'Novak Djokovic' AND m.market_type = '1x2';

-- Create market pools for all markets
INSERT INTO market_pools (market_id, total_pool, house_margin)
SELECT id, 0.00, margin_percentage FROM markets;
