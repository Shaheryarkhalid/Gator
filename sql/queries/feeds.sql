-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)VALUES($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name as feedName, feeds.url, users.name as userName from feeds join users on feeds.user_id = users.id GROUP BY feedName, url, userName order by userName asc;

-- name: GetFeedByUrl :one
SELECT id, name, url, user_id from feeds where url = $1 LIMIT 1;

-- name: MarkFeedFetched :exec

UPDATE feeds Set updated_at = Now(), last_fetched_at = Now () where id = $1;

-- name: DeleteFeedById :exec
DELETE  from feeds where id = $1;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;


-- name: GetNextFeedToFetch :one
SELECT * from feeds order by last_fetched_at asc NULLS FIRST LIMIT 1;
