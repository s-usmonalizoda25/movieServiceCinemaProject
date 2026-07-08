package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/db"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/logger"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/repository"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/server"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/service"
	moviepb "github.com/s-usmonalizoda25/movieServiceCinemaProject/moviepb/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	myLogger := logger.New()
	defer myLogger.Sync()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbPool, err := db.New(context.Background(), dsn)
	if err != nil {
		myLogger.Fatal("failed to connect to db and run migrations", zap.Error(err))
	}
	defer dbPool.Close()

	repo := repository.NewMovieRepository(dbPool)
	svc := service.New(repo, myLogger.Logger)
	movieServer := server.New(myLogger.Logger, svc.(*service.Service))

	grpcServer := grpc.NewServer()
	moviepb.RegisterMovieServiceServer(grpcServer, movieServer)
	reflection.Register(grpcServer)

	port := os.Getenv("SERVER_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		myLogger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		myLogger.Info("server started", zap.String("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			myLogger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	myLogger.Info("shutting down server...")

	_, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	grpcServer.GracefulStop()
	myLogger.Info("server stopped")
}
