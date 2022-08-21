package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/segmentio/kafka-go"

	"youtube-downloader-rest/internal/model"
)

func (tcg *TaskConsumerGroup) createTaskWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	r *kafka.Reader,
	w *kafka.Writer,
	wg *sync.WaitGroup,
	workerID int,
) {
	defer wg.Done()
	defer cancel()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			tcg.log.Errorf("FetchMessage", err)
			return
		}

		tcg.log.Infof(
			"WORKER: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n",
			workerID,
			m.Topic,
			m.Partition,
			m.Offset,
			string(m.Key),
			string(m.Value),
		)

		var task model.Task
		if err := json.Unmarshal(m.Value, &task); err != nil {
			tcg.log.Errorf("json.Unmarshall", err)
			continue
		}

		err = tcg.downloaderUC.UploadFile(ctx, &task)
		if err != nil {
			tcg.log.Errorf("taskUsecase.Create", err)
			continue
		}
		tcg.log.Infof("created task: %v", task)

		if err := r.CommitMessages(ctx, m); err != nil {
			tcg.log.Errorf("CommitMessages", err)
			continue
		}
	}
}
