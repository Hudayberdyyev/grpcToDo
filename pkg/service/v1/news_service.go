package v1

import (
	"context"

	v1 "github.com/Hudayberdyyev/grpcToDo/pkg/api/v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	apiVersion = "v1"
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
	err := s.db.Ping(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return nil, nil
}

func (s *newsServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	return &v1.ReadResponse{
		Api:      apiVersion,
		ExtraMsg: &v1.ExtraMsg{Id: 1, Title: "", Description: "", CreatedAt: &timestamppb.Timestamp{Seconds: 10, Nanos: 0}},
	}, nil
}

func (s *newsServiceServer) MustEmbedUnimplementedNewsServiceServer() {}
