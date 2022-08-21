package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	kafkaDownloader "youtube-downloader-rest/internal/downloader/controller/kafka"
	minioDownloader "youtube-downloader-rest/internal/downloader/controller/minio"
	"youtube-downloader-rest/internal/model"
	"youtube-downloader-rest/pkg/logger"
)

type downloaderUsecase struct {
	log         logger.Logger
	minioClient *minioDownloader.MinioClient
	dProducer   kafkaDownloader.DoneTaskProducer
}

func NewDownloderUsecase(log logger.Logger, minioClient *minioDownloader.MinioClient, dProducer kafkaDownloader.DoneTaskProducer) *downloaderUsecase {
	return &downloaderUsecase{
		log:         log,
		minioClient: minioClient,
		dProducer:   dProducer,
	}
}

func (d *downloaderUsecase) UploadFile(ctx context.Context, task *model.Task) error {
	taskBytes, err := json.Marshal(&task)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}

	_, err = d.minioClient.UploadFile(ctx, task)
	if err != nil {
		d.log.Errorf("minioClient.UploadFile", err)
		return err
	}

	return d.dProducer.PublishCreate(ctx, kafka.Message{
		Value: taskBytes,
		Time:  time.Now().UTC(),
	})
}
