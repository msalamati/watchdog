// Code generated by sqlc. DO NOT EDIT.
// source: deploys.sql

package deploys

import (
	"context"

	"github.com/jackc/pgtype"
)

const addDeploy = `-- name: AddDeploy :one
INSERT INTO deploys (
    service_name,
    status
) VALUES (
    $1,
    $2
)
RETURNING id, create_time, service_name, status
`

type AddDeployParams struct {
	ServiceName string
	Status      DeployStatus
}

func (q *Queries) AddDeploy(ctx context.Context, arg AddDeployParams) (Deploy, error) {
	row := q.db.QueryRowContext(ctx, addDeploy, arg.ServiceName, arg.Status)
	var i Deploy
	err := row.Scan(
		&i.ID,
		&i.CreateTime,
		&i.ServiceName,
		&i.Status,
	)
	return i, err
}

const getDeploy = `-- name: GetDeploy :one
SELECT id, create_time, service_name, status FROM deploys
WHERE id = $1
`

func (q *Queries) GetDeploy(ctx context.Context, id pgtype.UUID) (Deploy, error) {
	row := q.db.QueryRowContext(ctx, getDeploy, id)
	var i Deploy
	err := row.Scan(
		&i.ID,
		&i.CreateTime,
		&i.ServiceName,
		&i.Status,
	)
	return i, err
}
