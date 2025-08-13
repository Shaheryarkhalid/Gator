-- name: CreateFeedFollow :one
WITH inserted_feeds_follow as (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id , feed_id)VALUES($1, $2, $3, $4, $5) returning *
) SELECT inserted_feeds_follow.*, users.name as user_name, feeds.name as feed_name from inserted_feeds_follow join users on inserted_feeds_follow.user_id = users.id join feeds on inserted_feeds_follow.feed_id = feeds.id LIMIT 1; 

-- name: GetFeedFollowByFeedIdAndUserId :one
SELECT * from feed_follows where feed_id = $1 and user_id = $2 LIMIT 1;

-- name: GetFeedFollowsForUser :many
SELECT  feed_follows.*, users.name as user_name, feeds.name as feed_name from feed_follows  join users on feed_follows.user_id = users.id join feeds on feed_follows.feed_id = feeds.id where feed_follows.user_id = $1;

-- name: DeleteFeedFollowById :exec
DELETE from feed_follows where id = $1;

