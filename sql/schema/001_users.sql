-- +goose UP
CREATE TABLE users(id UUID Primary Key, created_at TIMESTAMP not null, updated_at TIMESTAMP not null, name TEXT not null unique);

-- +goose DOWN
DROP TABLE users;

