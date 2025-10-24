-- name: CreateFeedFollow :one
WITH new_follow AS (
    INSERT INTO feed_follows (
        id, created_at, updated_at, user_id, feed_id
    ) VALUES (
        DEFAULT,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        $1,
        $2
    )
    RETURNING *
)
SELECT
    nf.id,
    sqlc.embed(feeds),
    sqlc.embed(users)
FROM new_follow AS nf
JOIN feeds ON nf.feed_id = feeds.id
JOIN users ON nf.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT
    sqlc.embed(feeds), sqlc.embed(users)
FROM feed_follows AS follow
JOIN feeds ON follow.feed_id = feeds.id
JOIN users ON follow.user_id = users.id
WHERE users.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows USING feeds
WHERE feed_follows.user_id = $1 AND feed_id = feeds.id AND feeds.url = $2;
