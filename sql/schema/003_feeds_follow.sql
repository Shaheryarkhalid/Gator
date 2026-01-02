-- +goose UP
CREATE TABLE feeds_follow(
    id uuid primary key, 
    created_at timestamp not null, 
    updated_at timestamp not null, 
    user_id uuid not null, 
    feed_id uuid not null, 
    unique(user_id, feed_id), 
    constraint fk_user_id foreign key (user_id) references users(id) on delete cascade on update cascade,
    constraint fk_feed_id foreign key (feed_id) references feeds(id) on delete cascade on update cascade
);

-- +goose DOWN
DROP TABLE feeds_follow;

