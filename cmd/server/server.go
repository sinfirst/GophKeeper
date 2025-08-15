package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sinfirst/GophKeeper/internal/app"
	"github.com/sinfirst/GophKeeper/internal/config"
	"github.com/sinfirst/GophKeeper/internal/handlers"
	"github.com/sinfirst/GophKeeper/internal/middleware/logging"
	"github.com/sinfirst/GophKeeper/internal/storage"
	pb "github.com/sinfirst/GophKeeper/proto/gophkeeper"
	"google.golang.org/grpc"
)

func main() {
	config := config.NewConfig()
	logger := logging.NewLogger()
	storage := storage.NewPGDB(config, logger)
	handlers := handlers.NewHandler(storage, config, logger)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logging.LoggingUnaryInterceptor(logger)),
	)

	pb.RegisterGophKeeperServer(grpcServer, app.NewGophKeeperServer(handlers, logger))

	lis, err := net.Listen("tcp", config.Host)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	grpcServer.GracefulStop()
}
