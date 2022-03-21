-- name: GetChangesForUser :many
SELECT * FROM user_changes WHERE email = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetChangesForUserEarliestFirst :many
SELECT * FROM user_changes WHERE email = $1 ORDER BY created_at LIMIT $2 OFFSET $3;

-- name: GetChanges :many
SELECT * FROM user_changes ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetChangesEarliestFirst :many
SELECT * FROM user_changes ORDER BY created_at LIMIT $1 OFFSET $2;

-- name: GetChangesCount :one
SELECT COUNT(*) FROM user_changes;

-- name: GetChangesCountForUser :one
SELECT COUNT(*) FROM user_changes WHERE email = @email;

-- name: GetDeletions :many
SELECT * FROM user_deletions ORDER BY deleted_at DESC LIMIT $1 OFFSET $2;

-- name: GetDeletionsCount :one
SELECT COUNT(*) FROM user_deletions;

-- name: RegisterLoginAttempt :exec
INSERT INTO login_attempts (email, success, login_time) VALUES (@email, @success, now());

-- name: GetLoginAttemptsForUser :many
SELECT * FROM login_attempts
WHERE email = $1
ORDER BY login_time DESC
LIMIT $2 OFFSET $3;

-- name: GetFailedLoginAttemptsForUser :many
SELECT * FROM login_attempts
WHERE email = $1 AND success = FALSE
ORDER BY login_time DESC
LIMIT $2 OFFSET $3;