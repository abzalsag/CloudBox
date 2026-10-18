-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token_hash, user_agent, ip, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetActiveRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > now();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = now()
WHERE id = $1 AND revoked_at IS NULL;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: ListUserSessions :many
SELECT id, user_id, user_agent, ip, expires_at, created_at
FROM refresh_tokens
WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > now()
ORDER BY created_at DESC;

-- name: DeleteExpiredRefreshTokens :execrows
DELETE FROM refresh_tokens
WHERE expires_at < now() - interval '7 days' OR revoked_at < now() - interval '7 days';