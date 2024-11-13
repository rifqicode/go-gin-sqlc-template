-- name: CreateServerConfig :one
INSERT INTO server_configs (config_name, status, value, description)
VALUES ($1, $2, $3, $4)
RETURNING id, config_name, status, value, description, created_at, updated_at, deleted_at;

-- name: GetServerConfig :one
SELECT id, config_name, status, value, description, created_at, updated_at, deleted_at
FROM server_configs
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateServerConfig :one
UPDATE server_configs
SET config_name = $2, status = $3, value = $4, description = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, config_name, status, value, description, created_at, updated_at, deleted_at;

-- name: DeleteServerConfig :one
UPDATE server_configs
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, config_name, status, value, description, created_at, updated_at, deleted_at;

-- name: ListServerConfigs :many
SELECT id, config_name, status, value, description, created_at, updated_at, deleted_at
FROM server_configs
WHERE deleted_at IS NULL
ORDER BY created_at DESC;