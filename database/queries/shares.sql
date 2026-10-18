-- name: UpsertShare :one
INSERT INTO shares (id, resource_type, resource_id, owner_id, shared_with_user_id, permission)
VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (resource_type, resource_id, shared_with_user_id)
DO UPDATE SET permission = EXCLUDED.permission
           RETURNING *;

-- name: DeleteShare :exec
DELETE FROM shares WHERE id = $1 AND owner_id = $2;

-- name: ListSharesByResource :many
SELECT * FROM shares WHERE resource_type = $1 AND resource_id = $2 AND owner_id = $3;

-- name: ListSharedWithUser :many
SELECT * FROM shares WHERE shared_with_user_id = $1 ORDER BY created_at DESC;

-- name: GetUserPermission :one
SELECT permission FROM shares
WHERE resource_type = $1 AND resource_id = $2 AND shared_with_user_id = $3;