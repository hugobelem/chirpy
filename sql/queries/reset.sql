-- name: DeleteUsers :exec
DELETE FROM users
RETURNING id;