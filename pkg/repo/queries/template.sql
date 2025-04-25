-- name: CreateTemplate :one
INSERT INTO
    templates (id, name, repo_name, dockerfile)
VALUES
    (?, ?, ?, ?) RETURNING id;

-- name: DeleteTemplate :exec
DELETE FROM templates
WHERE
    id = ?;

-- name: UpdateTemplate :exec
UPDATE templates
SET
    name = ?,
    repo_name = ?,
    dockerfile = ?
WHERE
    id = ?;

-- name: GetTemplate :one
SELECT
    id,
    name,
    repo_name,
    dockerfile
FROM
    templates
WHERE
    id = ?;

-- name: GetAllTemplates :many
SELECT
    id,
    name,
    repo_name,
    dockerfile
FROM
    templates;