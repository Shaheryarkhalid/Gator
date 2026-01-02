-- name: CreateFeedFollow :one

with inserted_feed_follow AS (
    INSERT INTO feeds_follow(
        id, 
        created_at, 
        updated_at, 
        user_id, 
        feed_id
    ) values($1, $2, $3, $4, $5) returning *
) select 
        inserted_feed_follow.*, 
        users.name as user_name, 
        feeds.name as feed_name 
    from inserted_feed_follow 
    join 
        users on inserted_feed_follow.user_id = users.id
    join
        feeds on inserted_feed_follow.feed_id = feeds.id;



-- name: GetFeedFollowsForUser :many
select users.name as user_name, feeds.name as feed_name from feeds_follow  join users on feeds_follow.user_id = users.id join feeds on feeds_follow.feed_id = feeds.id where users.name = $1;

-- name: DeleteFeedFollow :exec
delete from feeds_follow where feed_id = $1 and user_id = $2;
