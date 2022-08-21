package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"youtube-downloader-rest/config"
	"youtube-downloader-rest/internal/rest"
	"youtube-downloader-rest/pb"
	"youtube-downloader-rest/pkg/logger"
)

const serviceName = "REST"

func main() {
	log.Printf("Starting %s microservice", serviceName)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Config
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	// gRPC
	conn, err := grpc.Dial("storage:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("fail to dial:", err)
	}
	defer conn.Close()
	client := pb.NewStorageClient(conn)

	// Http
	srv := rest.NewServer(client, appLogger, cfg)
	srv.Run()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
