package deploys

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	deployspb "github.com/msalamati/watchdog/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Directory struct {
	logger  *logrus.Logger
	db      *sql.DB
	sb      squirrel.StatementBuilderType
	querier Querier
}

func NewDirectory(logger *logrus.Logger, pgURL *url.URL) (*Directory, error) {
	connURL := *pgURL
	connURL.Scheme = "postgres"
	c, err := pgx.ParseConfig(connURL.String())
	if err != nil {
		return nil, fmt.Errorf("parsing postgres URI: %w", err)
	}

	c.Logger = logrusadapter.NewLogger(logger)
	db := stdlib.OpenDB(*c)

	err = validateSchema(db, pgURL.Scheme)
	if err != nil {
		return nil, fmt.Errorf("validating schema: %w", err)
	}

	return &Directory{
		logger:  logger,
		db:      db,
		sb:      squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		querier: New(db),
	}, nil
}

func (d Directory) Close() error {
	return d.db.Close()
}

func (d Directory) AddDeploy(ctx context.Context, req *deployspb.AddDeployRequest) (*deployspb.Deploy, error) {
	pgDeploy, err := d.querier.AddDeploy(ctx, AddDeployParams{
		ServiceName: req.GetServiceName(),
		Status: DeployStatusInProgress,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unexpected error adding deploy: %s", err.Error())
	}
	return deployPostgresToProto(pgDeploy)
}

func (d Directory) GetDeploy(ctx context.Context, req *deployspb.GetDeployRequest) (*deployspb.Deploy, error) {
	var id pgtype.UUID
	err := id.Set(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assign string to UUID: %s", err.Error())
	}

	pgDeploy, err := d.querier.GetDeploy(ctx, id)
	return deployPostgresToProto(pgDeploy)
}
