package downloader

import (
	"context"

	"youtube-downloader-rest/internal/model"
)

// Usecase describes the Storage service.
type Usecase interface {
	UploadFile(ctx context.Context, task *model.Task) error
}
