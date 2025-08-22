package app

import (
	"context"
	"errors"
	"fmt"

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

func (s *GophKeeperServer) Register(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, err := s.handlers.Register(ctx, req.Username, req.Password)
	if errors.Is(err, fmt.Errorf("conflict")) {
		return nil, status.Error(codes.AlreadyExists, "username already exist")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Server problem")
	}
	return &pb.AuthResponse{Token: token}, status.Error(codes.OK, "OK")

}

func (s *GophKeeperServer) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, err := s.handlers.Login(ctx, req.Username, req.Password)
	if errors.Is(err, fmt.Errorf("unauthenticated")) {
		return nil, status.Error(codes.Unauthenticated, "uncorrect login or password")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Server problem")
	}
	return &pb.AuthResponse{Token: token}, status.Error(codes.OK, "OK")
}

func (s *GophKeeperServer) StoreData(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	id, err := s.handlers.StoreData(ctx, req.Token, models.Record{TypeRecord: req.Record.Type, Data: req.Record.Data, Meta: req.Record.Meta})
	if errors.Is(err, fmt.Errorf("unauthenticated")) {
		return nil, status.Error(codes.Unauthenticated, "no valid token")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Server problem")
	}
	return &pb.StoreResponse{Id: int64(id)}, status.Error(codes.OK, "OK")
}

func (s *GophKeeperServer) RetrieveData(ctx context.Context, req *pb.RetrieveRequest) (*pb.RetrieveResponse, error) {
	return nil, nil
}
func (s *GophKeeperServer) ListData(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	return nil, nil
}
func (s *GophKeeperServer) DeleteData(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
