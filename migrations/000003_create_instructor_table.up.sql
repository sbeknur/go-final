CREATE TABLE IF NOT EXISTS instructors (
    id bigserial PRIMARY KEY,
    firstName text NOT NULL,
    lastName text NOT NULL,
    age integer NOT NULL
);
