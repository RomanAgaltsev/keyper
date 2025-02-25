package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	recordsv1 "github.com/RomanAgaltsev/keyper/pkg/records_v1"
	"github.com/RomanAgaltsev/keyper/server/internal/converter"
	"github.com/RomanAgaltsev/keyper/server/internal/model"
)

type Records interface {
	Get(ctx context.Context, id uint64) (model.Record, error)
	List(ctx context.Context, login string, recType model.RecordType) ([]model.Record, error)
}

func Register(gRPCServer *grpc.Server, records Records) {
	recordsv1.RegisterRecordsV1Server(gRPCServer, &gRPCServerAPI{records: records})
}

type gRPCServerAPI struct {
	recordsv1.UnimplementedRecordsV1Server
	records Records
}

func (s *gRPCServerAPI) Get(ctx context.Context, in *recordsv1.GetRequest) (*recordsv1.GetResponse, error) {
	if in.RecordId == 0 {
		return nil, status.Error(codes.InvalidArgument, "record ID is required")
	}

	record, err := s.records.Get(ctx, in.GetRecordId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get record")
	}

	return &recordsv1.GetResponse{
		Record: converter.ToRecordsFromServiceRecord(&record),
	}, nil
}

func (s *gRPCServerAPI) List(in *recordsv1.ListRequest, server recordsv1.RecordsV1_ListServer) error {
	if in.Login == "" {
		return status.Error(codes.InvalidArgument, "user login is required")
	}

	records, err := s.records.List(context.Background(), in.GetLogin(), model.RecordType(in.GetType()))
	if err != nil {
		return status.Error(codes.Internal, "failed to get records")
	}

	for _, record := range records {
		if err := server.Send(&recordsv1.ListResponse{
			Record: converter.ToRecordsFromServiceRecord(&record),
		}); err != nil {
			return status.Error(codes.Internal, "failed to send record")
		}
	}
	return nil
}
