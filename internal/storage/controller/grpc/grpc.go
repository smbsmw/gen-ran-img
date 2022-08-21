package grpc

import (
	"context"

	"youtube-downloader-rest/internal/model"
	"youtube-downloader-rest/internal/storage"
	"youtube-downloader-rest/pb"
	"youtube-downloader-rest/pkg/logger"
)

type storageService struct {
	log    logger.Logger
	taskUC storage.Usecase
	pb.UnimplementedStorageServer
}

func NewTaskService(log logger.Logger, taskUC storage.Usecase) *storageService {
	return &storageService{
		log:    log,
		taskUC: taskUC,
	}
}

func (s *storageService) CreateTask(ctx context.Context, task *pb.Task) (*pb.Feature, error) {
	s.log.Info("Get info")
	t, err := model.TaskFromProto(task)
	if err != nil {
		s.log.Error("TaskFromProto", err)
		return nil, err
	}
	err = s.taskUC.Create(ctx, t)
	if err != nil {
		s.log.Error("taskUC.Create", err)
	}

	return &pb.Feature{Name: task.ImageName}, err
}
