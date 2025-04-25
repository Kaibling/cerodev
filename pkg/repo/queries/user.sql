-- name: CreateUser :one
INSERT INTO
	users (id, username, password)
VALUES
	(?, ?, ?) RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
	id = ?;

-- name: UpdateUser :exec
UPDATE users
SET
	username = ?,
	password = ?
WHERE
	id = ?;

-- name: GetUserByID :many
SELECT
	u.id,
	u.username,
	t.token
FROM
	users u
	LEFT JOIN tokens t ON u.id = t.user_id
WHERE
	id = ?;

-- name: GetUnsafeUserByUsername :many
SELECT
	u.id,
	u.username,
	u.password,
	t.token
FROM
	users u
	LEFT JOIN tokens t ON u.id = t.user_id
WHERE
	username = ?;

-- name: GetAllUsers :many
SELECT
	u.id,
	u.username,
	t.token
FROM
	users u
	LEFT JOIN tokens t ON u.id = t.user_id;

-- name: GetUserByToken :many
SELECT
	u.id,
	u.username,
	t.token
FROM
	tokens t
	JOIN users u ON t.user_id = u.id
WHERE
	u.id = (
		SELECT
			user_id
		FROM
			tokens
		WHERE
			tokens.token = ?
	);