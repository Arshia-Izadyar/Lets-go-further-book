ALTER TABLE movies DROP CONSTRAINT IF EXISTS  movie_runtime_check;
ALTER TABLE movies DROP CONSTRAINT IF EXISTS  movie_year_check;
ALTER TABLE movies DROP CONSTRAINT IF EXISTS  genre_length_check;