-- +goose UP

CREATE TABLE posts(id uuid Primary key, created_at timestamp not null, updated_at timestamp not null, title text not null, url text not null, description text not null, published_at timestamp not null, feed_id uuid not null, unique(title, url), constraint fk_feed_id foreign key (feed_id) references feeds(id));

-- +goose DOWN
DROP TABLE posts;
