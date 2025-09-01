package app

import (
	"context"
	"errors"

	"github.com/sinfirst/GophKeeper/internal/config"
	"github.com/sinfirst/GophKeeper/internal/handlers"
	"github.com/sinfirst/GophKeeper/internal/models"
	pb "github.com/sinfirst/GophKeeper/proto/gophkeeper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GophKeeperServer struct {
	pb.UnimplementedGophKeeperServer
	logger   zap.SugaredLogger
	handlers handlers.Handler
}

func NewGophKeeperServer(handlers handlers.Handler, logger zap.SugaredLogger) pb.GophKeeperServer {
	return &GophKeeperServer{handlers: handlers, logger: logger}
}

func (s *GophKeeperServer) GetVersion(ctx context.Context, req *emptypb.Empty) (*pb.GetVersionResponse, error) {
	return &pb.GetVersionResponse{Ver: &pb.Version{Version: config.VersionBuild, Date: config.DateBuild}}, status.Error(codes.OK, "OK")
}
func (s *GophKeeperServer) Register(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, err := s.handlers.Register(ctx, req.Username, req.Password)
	if err = s.errorHandler(err); err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Token: token}, status.Error(codes.OK, "OK")

}

func (s *GophKeeperServer) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, err := s.handlers.Login(ctx, req.Username, req.Password)
	if err = s.errorHandler(err); err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Token: token}, status.Error(codes.OK, "OK")
}

func (s *GophKeeperServer) StoreData(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	id, err := s.handlers.StoreData(ctx, req.Token, models.Record{TypeRecord: req.Record.Type, Data: req.Record.Data, Meta: req.Record.Meta})
	if err = s.errorHandler(err); err != nil {
		return nil, err
	}

	return &pb.StoreResponse{Id: int64(id)}, status.Error(codes.OK, "OK")
}

func (s *GophKeeperServer) RetrieveData(ctx context.Context, req *pb.RetrieveRequest) (*pb.RetrieveResponse, error) {
	record, err := s.handlers.RetrieveData(ctx, req.Token, int(req.Id))
	if err = s.errorHandler(err); err != nil {
		return nil, err
	}

	return &pb.RetrieveResponse{Record: &pb.DataRecord{Type: record.TypeRecord, Data: record.Data, Meta: record.Meta}}, status.Error(codes.OK, "OK")
}
func (s *GophKeeperServer) UpdateData(ctx context.Context, req *pb.UpdateResponse) (*emptypb.Empty, error) {
	err := s.handlers.UpdateData(ctx, req.Token, req.Meta, int(req.Id), req.Data)
	if err = s.errorHandler(err); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.OK, "OK")
}
func (s *GophKeeperServer) ListData(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	records, err := s.handlers.ListData(ctx, req.Token)
	if err = s.errorHandler(err); err != nil {
		return nil, err
	}
	if records == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	var resp []*pb.DataRecord
	for _, i := range records {
		resp = append(resp, &pb.DataRecord{
			Id:   int64(i.Id),
			Type: i.TypeRecord,
			Data: i.Data,
			Meta: i.Meta,
		})
	}
	return &pb.ListResponse{Records: resp}, status.Error(codes.OK, "OK")

}
func (s *GophKeeperServer) DeleteData(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := s.handlers.DeleteData(ctx, req.Token, int(req.Id))
	if err = s.errorHandler(err); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.OK, "OK")

}

func (s *GophKeeperServer) errorHandler(err error) error {
	var appErr models.AppError
	if errors.As(err, &appErr) {
		switch appErr {
		case models.ErrUnauthenticated:
			return status.Error(codes.Unauthenticated, string(appErr))
		case models.ErrConflict:
			return status.Error(codes.AlreadyExists, string(appErr))
		case models.ErrAccessDenied:
			return status.Error(codes.PermissionDenied, string(appErr))
		case models.ErrNotFound:
			return status.Error(codes.NotFound, string(appErr))
		}
	}

	if err != nil {
		s.logger.Errorf("err: %v", err)
		return status.Error(codes.Internal, "Server problem")
	}
	return nil
}
