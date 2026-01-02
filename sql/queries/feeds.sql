-- name: CreateFeed :one

INSERT INTO feeds(id, created_at, updated_at, name, url, user_id ) values($1, $2, $3, $4,$5, $6) returning *;

-- name: GetFeedsByUrl :one
SELECT * from feeds  where url = $1 limit 1;

-- name: GetFeedsByUsers :many
SELECT feeds.name, feeds.url, users.name as user_name from feeds join users on feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
update feeds set last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP where id = $1;

-- name: GetNextFeedToFetch :one
SELECT * from feeds order by last_fetched_at  asc NULLS FIRST limit 1;

