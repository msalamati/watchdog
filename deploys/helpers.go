package deploys

import (
	"database/sql"
	"embed"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	deployspb "github.com/msalamati/watchdog/proto"
)

//go:embed migrations/*.sql
var fs embed.FS

const version = 1

func validateSchema(db *sql.DB, scheme string) error {
	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	var driverInstance database.Driver
	driverInstance, err = postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceInstance, scheme, driverInstance)
	if err != nil {
		return err
	}
	err = m.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return sourceInstance.Close()
}

func deployPostgresToProto(pgDeploy Deploy) (*deployspb.Deploy, error) {
	protoStatus, err := deployStatusPostgresToProto(pgDeploy.Status)
	if err != nil {
		return nil, err
	}
	var userID string
	err = pgDeploy.ID.AssignTo(&userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assign UUID to string: %s", err.Error())
	}
	return &deployspb.Deploy{
		CreateTime:  timestamppb.New(pgDeploy.CreateTime),
		Id:          userID,
		ServiceName: pgDeploy.ServiceName,
		Status:      protoStatus,
	}, nil
}

func deployStatusPostgresToProto(pgStatus DeployStatus) (deployspb.DeployStatus, error) {
	switch pgStatus {
	case DeployStatusInProgress:
		return deployspb.DeployStatus_IN_PROGRESS, nil
	case DeployStatusSucceeded:
		return deployspb.DeployStatus_SUCCEEDED, nil
	case DeployStatusFailed:
		return deployspb.DeployStatus_FAILED, nil
	default:
		return 0, status.Errorf(codes.Internal, "unknown status type %q", pgStatus)
	}
}
