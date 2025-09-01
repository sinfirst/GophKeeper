-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    username TEXT NOT NULL PRIMARY KEY,
    user_password TEXT NOT NULL
);

CREATE TYPE record_type AS ENUM (
    'LOGIN',
    'TEXT',
    'CARD',
    'BINARY'
);

CREATE TABLE IF NOT EXISTS records (
    id SERIAL NOT NULL PRIMARY KEY,
    type_record record_type NOT NULL,
    user_data BYTEA NOT NULL,
    meta TEXT,
    username TEXT REFERENCES users(username) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS record_type;
DROP TABLE IF EXISTS records;
-- +goose StatementEnd