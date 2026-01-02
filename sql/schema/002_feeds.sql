-- +goose UP

CREATE TABLE feeds(id uuid Primary Key, created_at timestamp not null, updated_at timestamp not null, name text not null,  url text not null unique, user_id uuid not null, constraint fk_user_id Foreign key (user_id) references users(id) on delete cascade on update cascade);
-- +goose DOWN
Drop TABLE feeds;
