package storage

import (
	"context"

	"youtube-downloader-rest/internal/model"
)

// Usecase describes the Storage service.
type Usecase interface {
	Create(ctx context.Context, url *model.Task) error
	CheckTaskStatus(ctx context.Context, url *model.Task) (bool, error)
}
