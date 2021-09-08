package v1

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/Hudayberdyyev/grpcToDo/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
	dbURL      = "postges://postgres:@%!)(^@localhost:5432/postgres?sslmode=disable"
)

type newsServiceServer struct {
	db *pgxpool.Pool
}

func NewNewsServiceServer(db *pgxpool.Pool) v1.NewsServiceServer {
	return &newsServiceServer{db: db}
}

func (s *newsServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func (s *newsServiceServer) connect(ctx context.Context) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return db, nil
}

func (s *newsServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query ToDo by ID
	rows, err := c.Query(ctx, "SELECT `ID`, `Title`, `Description`, `Created_at` FROM extra_messages WHERE `ID`=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from news-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from news-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("msg with ID='%d' is not found",
			req.Id))
	}

	// get ToDo data
	var td v1.ExtraMsg
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
	}
	td.CreatedAt, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple news rows with ID='%d'",
			req.Id))
	}

	return &v1.ReadResponse{
		Api:      apiVersion,
		ExtraMsg: &td,
	}, nil
}

func (s *newsServiceServer) MustEmbedUnimplementedNewsServiceServer() {}
