CREATE INDEX IF NOT EXISTS courses_title_idx ON courses USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS courses_lectures_idx ON courses USING GIN (lectures);
