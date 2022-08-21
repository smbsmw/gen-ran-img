package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"youtube-downloader-rest/internal/model"
	kafkaStorage "youtube-downloader-rest/internal/storage/controller/kafka"
	"youtube-downloader-rest/pkg/logger"
)

type storageUsecase struct {
	log       logger.Logger
	sProducer kafkaStorage.TaskProducer
}

func NewTaskUC(log logger.Logger, sProducer kafkaStorage.TaskProducer) *storageUsecase {
	return &storageUsecase{
		log:       log,
		sProducer: sProducer,
	}
}

func (s *storageUsecase) Create(ctx context.Context, task *model.Task) error {
	taskBytes, err := json.Marshal(&task)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
		//return errors.Wrap(err, "json.Marshal")
	}

	return s.sProducer.PublishCreate(ctx, kafka.Message{
		Value: taskBytes,
		Time:  time.Now().UTC(),
	})
}

func (s *storageUsecase) CheckTaskStatus(ctx context.Context, task *model.Task) (bool, error) {
	return false, nil
}
