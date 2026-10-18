-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, full_name)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE lower(email) = lower($1);

-- name: UpdateUserProfile :one
UPDATE users
SET full_name = $2, updated_at = now()
WHERE id = $1
    RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1;

-- name: IncreaseStorageUsed :one

UPDATE users
SET storage_used = storage_used + sqlc.arg(delta)::bigint, updated_at = now()
WHERE id = sqlc.arg(id)
  AND storage_used + sqlc.arg(delta)::bigint <= storage_quota
    RETURNING storage_used, storage_quota;

-- name: DecreaseStorageUsed :exec
UPDATE users
SET storage_used = GREATEST(storage_used - sqlc.arg(delta)::bigint, 0), updated_at = now()
WHERE id = sqlc.arg(id);