-- name: CreateFile :one
INSERT INTO files (
    user_id,
    name,
    storage_key,
    size,
    content_type
) VALUES (
             $1,
             $2,
             $3,
             $4,
             $5
         )
    RETURNING id, user_id, name, storage_key, size, content_type, created_at, updated_at;

-- name: GetFileByID :one
SELECT id, user_id, name, storage_key, size, content_type, created_at, updated_at
FROM files
WHERE id = $1;

-- name: ListFilesByUserID :many
SELECT id, user_id, name, storage_key, size, content_type, created_at, updated_at
FROM files
WHERE user_id = $1
ORDER BY id DESC;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1 AND user_id = $2;