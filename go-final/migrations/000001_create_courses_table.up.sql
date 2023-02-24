CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS courses (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    published_date text NOT NULL,
    runtime integer NOT NULL,
    lectures text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);
