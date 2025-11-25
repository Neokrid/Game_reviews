CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_games_title_trgm ON games USING gin (title gin_trgm_ops);