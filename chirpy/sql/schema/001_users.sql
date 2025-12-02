-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email varchar(400) NOT NULL,
    UNIQUE (email)
);

-- +goose Down
DROP TABLE users;
