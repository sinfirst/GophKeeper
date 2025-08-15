-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    user_id INT PRIMARY KEY UNIQUE,
    username TEXT NOT NULL UNIQUE,
    user_password TEXT NOT NULL UNIQUE,
    token TEXT UNIQUE
);

CREATE TYPE record_type AS ENUM (
    'LOGIN',
    'TEXT',
    'CARD',
    'BINARY'
);

CREATE TABLE IF NOT EXISTS records (
    id INT PRIMARY KEY NOT NULL UNIQUE,
    user_id INT NOT NULL,
    type_record record_type NOT NULL,
    user_data BYTEA NOT NULL,
    meta TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS records ;
-- +goose StatementEnd