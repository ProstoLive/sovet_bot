CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    first_name TEXT,
    last_name TEXT,
    username TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE user_applications (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    answers JSONB,
    submitted_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE application_states (
    user_id BIGINT PRIMARY KEY,
    current_question INTEGER, 
    answers JSONB DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE events {
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    place TEXT,
    datetime TIMESTAMP,
    author_id BIGINT
}