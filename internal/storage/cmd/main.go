package main

import (
	"context"
	"net"
	"os"
	"log"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"youtube-downloader-rest/config"
	kafkaPkg "youtube-downloader-rest/pkg/kafka"
	grpcStorage "youtube-downloader-rest/internal/storage/controller/grpc"
	"youtube-downloader-rest/internal/storage/controller/kafka"
	"youtube-downloader-rest/internal/storage/usecase"
	"youtube-downloader-rest/pb"
	"youtube-downloader-rest/pkg/logger"
)

const port = ":50051"

func main() {
	log.Println("Starting storage microservice")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Config
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	// Kafka
	conn, err := kafkaPkg.NewKafkaConn(cfg)
	if err != nil {
		log.Fatal("NewKafkaConn", err)
	}
	defer conn.Close()
	brokers, err := conn.Brokers()
	if err != nil {
		log.Fatal("conn.Brokers", err)
	}
	log.Printf("Kafka connected: %v", brokers)

	sProducer := kafka.NewTaskProducer(appLogger, cfg)
	sProducer.Run()

	// Usecase
	taskUC := usecase.NewTaskUC(appLogger, sProducer)

	// Kafka
	sConsumer := kafka.NewTaskConsumerGroup(cfg.Kafka.Brokers, "StorageGroup", appLogger, cfg, taskUC)
	sConsumer.RunConsumers(ctx, cancel)

	// gRPC
	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	taskService := grpcStorage.NewTaskService(appLogger, taskUC)
	pb.RegisterStorageServer(srv, taskService)

	go func() {
		listener, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to create listener: %s",
				err)
		}
		log.Println("start server on port", port)
		if err := srv.Serve(listener); err != nil {
			log.Println("failed to exit serve: ", err)
		}
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	// Shutdown
	srv.GracefulStop()
}
