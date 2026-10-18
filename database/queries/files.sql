-- name: CreateFile :one
INSERT INTO files (id, owner_id, folder_id, name, size, mime_type, checksum, storage_key, bucket)
VALUES ($1, $2, sqlc.narg('folder_id'), $3, $4, $5, $6, $7, $8)
    RETURNING *;

-- name: GetFile :one
SELECT * FROM files
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL;

-- name: GetFileIncludingDeleted :one
SELECT * FROM files WHERE id = $1 AND owner_id = $2;

-- name: ListFilesByFolder :many
SELECT * FROM files
WHERE owner_id = $1
  AND folder_id IS NOT DISTINCT FROM sqlc.narg('folder_id')
    AND deleted_at IS NULL
ORDER BY name
    LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: SearchFilesByName :many
SELECT * FROM files
WHERE owner_id = $1 AND deleted_at IS NULL AND name ILIKE '%' || sqlc.arg('query')::text || '%'
ORDER BY updated_at DESC
    LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: RenameFile :one
UPDATE files SET name = $3, updated_at = now()
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
    RETURNING *;

-- name: MoveFile :one
UPDATE files SET folder_id = sqlc.narg('folder_id'), updated_at = now()
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
    RETURNING *;

-- name: ReplaceFileContent :one
-- Жаңа нұсқа жүктелгенде (versioning болса, ескісін алдымен file_versions-қа көшіріңіз)
UPDATE files
SET size = $3, checksum = $4, storage_key = $5, mime_type = $6, updated_at = now()
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
    RETURNING *;

-- name: SoftDeleteFile :one
UPDATE files SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
    RETURNING *;

-- name: SoftDeleteFilesInFolders :exec
-- Папка өшірілгенде ішіндегі файлдар да Trash-қа кетеді
UPDATE files SET deleted_at = now(), updated_at = now()
WHERE owner_id = $1 AND folder_id = ANY(sqlc.arg('folder_ids')::uuid[]) AND deleted_at IS NULL;

-- name: RestoreFile :one
UPDATE files SET deleted_at = NULL, updated_at = now()
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NOT NULL
    RETURNING *;

-- name: RestoreFilesInFolders :exec
UPDATE files SET deleted_at = NULL, updated_at = now()
WHERE owner_id = $1 AND folder_id = ANY(sqlc.arg('folder_ids')::uuid[]);

-- name: ListTrashedFiles :many
-- Папкасы өшірілмеген, өзі өшірілген файлдар
SELECT f.* FROM files f
WHERE f.owner_id = $1 AND f.deleted_at IS NOT NULL
  AND (f.folder_id IS NULL
    OR NOT EXISTS (SELECT 1 FROM folders d WHERE d.id = f.folder_id AND d.deleted_at IS NOT NULL))
ORDER BY f.deleted_at DESC;

-- name: HardDeleteFile :exec
DELETE FROM files WHERE id = $1 AND owner_id = $2;

-- name: ListFilesInFolders :many
SELECT id, owner_id, size, bucket, storage_key FROM files
WHERE owner_id = $1 AND folder_id = ANY(sqlc.arg('folder_ids')::uuid[]);

-- name: ListExpiredTrashedFiles :many
SELECT id, owner_id, size, bucket, storage_key FROM files
WHERE deleted_at IS NOT NULL AND deleted_at < now() - interval '30 days'
    LIMIT $1;

-- name: CountFilesByStorageKey :one

SELECT count(*) FROM files WHERE bucket = $1 AND storage_key = $2;

-- name: FindFileByChecksum :one
SELECT * FROM files WHERE checksum = $1 AND size = $2 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserFolderStats :one
SELECT count(*)::bigint AS file_count, COALESCE(sum(size), 0)::bigint AS total_size
FROM files WHERE owner_id = $1 AND deleted_at IS NULL;