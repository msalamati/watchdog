-- name: AddDeploy :one
INSERT INTO deploys (
    service_name,
    status
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetDeploy :one
SELECT * FROM deploys
WHERE id = $1;
