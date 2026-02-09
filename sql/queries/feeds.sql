-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, api_key)
VALUES($1, $2, $3, $4,
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;
