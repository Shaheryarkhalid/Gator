-- name: CreateUser :one
INSERT INTO users(
    id, 
    created_at, 
    updated_at, 
    name
)
values(
    $1, 
    $2, 
    $3, 
    $4
)
returning  *;



-- name: GetUser :one
SELECT * FROM users where name = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: DeleteUsers :exec
DELETE  FROM users;
