package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"youtube-downloader-rest/config"
	kafkaDownloader "youtube-downloader-rest/internal/downloader/controller/kafka"
	minioDownloader "youtube-downloader-rest/internal/downloader/controller/minio"
	"youtube-downloader-rest/internal/downloader/usecase"
	kafkaPkg "youtube-downloader-rest/pkg/kafka"
	"youtube-downloader-rest/pkg/logger"
)

const port = ":50052"

func main() {
	log.Println("Starting Downloader microservice")

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

	// Minio
	minioClient := minioDownloader.NewMinioClient(cfg.Minio.Addr, "minioadmin", "minioadmin", false)

	// Kafka
	conn, err := kafkaPkg.NewKafkaConn(cfg)
	if err != nil {
		appLogger.Fatalf("NewKafkaConn", err)
	}
	defer conn.Close()
	brokers, err := conn.Brokers()
	if err != nil {
		appLogger.Fatalf("conn.Brokers", err)
	}
	appLogger.Infof("Kafka connected: %v", brokers)

	dProducer := kafkaDownloader.NewDoneTaskProducer(appLogger, cfg)
	dProducer.Run()

	// Usecase
	downloaderUC := usecase.NewDownloderUsecase(appLogger, minioClient, dProducer)

	//Kafka
	dConsumer := kafkaDownloader.NewTaskConsumerGroup(cfg.Kafka.Brokers, "DownloaderGroup", appLogger, cfg, downloaderUC)
	dConsumer.RunConsumers(ctx, cancel)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
