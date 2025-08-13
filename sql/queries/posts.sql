-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title,url, description , published_at , feed_id)VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *;

-- name: GetPostsForUser :many
SELECT * from posts where feed_id in (SELECT feed_id from feed_follows where user_id = $1) limit $2;
