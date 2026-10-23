-- name: CreateFolder :one
INSERT INTO folders (id, owner_id, parent_id, name)
VALUES ($1, $2, sqlc.narg('parent_id'), $3)
    RETURNING *;

-- name: GetFolder :one
SELECT * FROM folders
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL;

-- name: GetFolderIncludingDeleted :one
SELECT * FROM folders WHERE id = $1 AND owner_id = $2;

-- name: ListFoldersByParent :many

SELECT * FROM folders
WHERE owner_id = $1
  AND parent_id IS NOT DISTINCT FROM sqlc.narg('parent_id')
    AND deleted_at IS NULL
ORDER BY name
    LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: RenameFolder :one
UPDATE folders SET name = $3, updated_at = now()
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
    RETURNING *;

-- name: MoveFolder :one

UPDATE folders SET parent_id = sqlc.narg('parent_id'), updated_at = now()
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
    RETURNING *;

--- name: IsFolderDescendant :one
WITH RECURSIVE tree AS (
    SELECT f.id, f.parent_id FROM folders f WHERE f.id = sqlc.arg('candidate_id')
    UNION ALL
    SELECT f.id, f.parent_id FROM folders f JOIN tree t ON f.id = t.parent_id
)
SELECT EXISTS (
    SELECT 1 FROM tree WHERE tree.id = sqlc.arg('ancestor_id')
)::boolean AS is_descendant;

-- name: GetFolderBreadcrumbs :many

WITH RECURSIVE path AS (
    SELECT id, parent_id, name, 1 AS depth
    FROM folders WHERE folders.id = $1 AND folders.owner_id = $2
    UNION ALL
    SELECT f.id, f.parent_id, f.name, p.depth + 1
    FROM folders f JOIN path p ON f.id = p.parent_id
)
SELECT id, parent_id, name FROM path ORDER BY depth DESC;

-- name: ListFolderTreeIDs :many

WITH RECURSIVE tree AS (
    SELECT id FROM folders WHERE folders.id = $1 AND folders.owner_id = $2
    UNION ALL
    SELECT f.id FROM folders f JOIN tree t ON f.parent_id = t.id
)
SELECT id FROM tree;

-- name: SoftDeleteFolders :exec
UPDATE folders SET deleted_at = now(), updated_at = now()
WHERE id = ANY(sqlc.arg('ids')::uuid[]) AND owner_id = $1 AND deleted_at IS NULL;

-- name: RestoreFolders :exec
UPDATE folders SET deleted_at = NULL, updated_at = now()
WHERE id = ANY(sqlc.arg('ids')::uuid[]) AND owner_id = $1;

-- name: ListTrashedRootFolders :many

SELECT f.* FROM folders f
WHERE f.owner_id = $1 AND f.deleted_at IS NOT NULL
  AND (f.parent_id IS NULL
    OR NOT EXISTS (SELECT 1 FROM folders p WHERE p.id = f.parent_id AND p.deleted_at IS NOT NULL))
ORDER BY f.deleted_at DESC;

-- name: HardDeleteFolders :exec
DELETE FROM folders WHERE id = ANY(sqlc.arg('ids')::uuid[]) AND owner_id = $1;

-- name: ListExpiredTrashedFolderIDs :many

SELECT id, owner_id FROM folders
WHERE deleted_at IS NOT NULL AND deleted_at < now() - interval '30 days'
    LIMIT $1;