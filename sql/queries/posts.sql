-- name: CreatePost :one

 INSERT INTO posts(id, created_at , updated_at , title , url , description , published_at , feed_id) values ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2, $3, $4, $5, $6) returning *;

-- name: GetPostsForUser :many

SELECT posts.* from posts join feeds on posts.feed_id = feeds.id where feeds.user_id = $1 LIMIT $2;
