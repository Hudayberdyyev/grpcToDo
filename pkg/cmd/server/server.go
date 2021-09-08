package server

import (
	"context"
	"flag"
	"fmt"
	"log"

	// mysql driver
	v1 "github.com/Hudayberdyyev/grpcToDo/pkg/service/v1"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Hudayberdyyev/grpcToDo/pkg/protocol/grpc"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string
	// DatastoreDBSchema is schema of database
	DatastoreDBSchema string
	// DatastoreDBDriver is driver of database
	DatastoreDBDriver string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBDriver, "db-driver", "", "Database driver")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "sslmode=disable"

	dsn := fmt.Sprintf("%s://%s:%s@%s/%s?%s",
		cfg.DatastoreDBDriver,
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)
	log.Printf("%s\n", dsn)
	db, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewNewsServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
