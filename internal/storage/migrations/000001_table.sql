-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    username TEXT PRIMARY KEY NOT NULL UNIQUE,
    user_password TEXT NOT NULL UNIQUE,
);

CREATE TYPE record_type AS ENUM (
    'LOGIN',
    'TEXT',
    'CARD',
    'BINARY'
);

CREATE TABLE IF NOT EXISTS records (
    id INT PRIMARY KEY NOT NULL UNIQUE,
    type_record record_type NOT NULL,
    user_data BYTEA NOT NULL,
    meta TEXT
    username TEXT REFERENCES users(username) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS records;
-- +goose StatementEnd