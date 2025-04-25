-- name: CreateContainer :one
INSERT INTO
    containers (
        id,
        docker_id,
        image_name,
        container_name,
        git_repo,
        user_id,
        env_vars,
        ports
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?) returning id;

-- name: DeleteContainer :exec
DELETE FROM containers
WHERE
    id = ?;

-- name: GetContainerByID :one
SELECT
    c.id,
    c.docker_id,
    c.image_name,
    c.container_name,
    c.git_repo,
    c.user_id,
    c.env_vars,
    c.ports,
    p.port as ui_port
FROM
    containers c
    JOIN ports p on p.container_id = c.id
WHERE
    id = ?;

-- name: GetAllContainers :many
SELECT
    c.id,
    c.docker_id,
    c.image_name,
    c.container_name,
    c.git_repo,
    c.user_id,
    c.env_vars,
    c.ports,
    p.port as ui_port
FROM
    containers c
    JOIN ports p on p.container_id = c.id;

-- name: GetPortByContainerID :one
SELECT
    port
FROM
    ports
WHERE
    container_id = ?;

-- name: GetFreePort :one
SELECT
    port
FROM
    ports
WHERE
    in_use = 0
ORDER BY
    port
LIMIT
    1;

-- name: AllocatePort :exec
UPDATE ports
SET
    in_use = 1,
    container_id = ?
WHERE
    port = ?;

-- name: ReleasePortbyPort :exec
UPDATE ports
SET
    in_use = 0,
    container_id = NULL
WHERE
    port = ?;

-- name: ReleasePortByContainer :exec
UPDATE ports
SET
    in_use = 0,
    container_id = NULL
WHERE
    container_id = ?;

-- name: UpdateContainer :exec
UPDATE containers
SET
    docker_id = ?,
    image_name = ?,
    container_name = ?,
    git_repo = ?,
    user_id = ?,
    env_vars = ?,
    ports = ?
WHERE
    id = ?;

-- name: GetPortCount :one
SELECT
    count(port)
FROM
    ports;

-- name: CreatePort :exec
INSERT INTO
    ports (port, in_use, container_id)
VALUES
    (?, 0, NULL);