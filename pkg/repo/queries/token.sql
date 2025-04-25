-- name: CreateToken :exec
INSERT INTO
	tokens (user_id, token)
VALUES
	(?, ?);

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE
	token = ?;

-- name: GetToken :one
SELECT
	user_id,
	token
FROM
	tokens
WHERE
	token = ?;

-- name: GetTokenByUserID :many
SELECT
	user_id,
	token
FROM
	tokens
WHERE
	user_id = ?;