package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/segmentio/kafka-go"

	"youtube-downloader-rest/config"
	"youtube-downloader-rest/internal/model"
	"youtube-downloader-rest/internal/storage"
	"youtube-downloader-rest/pkg/logger"
)

type TaskConsumerGroup struct {
	Brokers     []string
	GroupID     string
	log         logger.Logger
	cfg         *config.Config
	taskUsecase storage.Usecase
}

func NewTaskConsumerGroup(
	brokers []string,
	groupID string,
	log logger.Logger,
	cfg *config.Config,
	taskUsecase storage.Usecase,
) *TaskConsumerGroup {
	return &TaskConsumerGroup{
		Brokers:     brokers,
		GroupID:     groupID,
		log:         log,
		cfg:         cfg,
		taskUsecase: taskUsecase,
	}
}

func (tcg *TaskConsumerGroup) getNewKafkaReader(kafkaURL []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		Logger:                 kafka.LoggerFunc(tcg.log.Debugf),
		ErrorLogger:            kafka.LoggerFunc(tcg.log.Errorf),
		MaxAttempts:            maxAttempts,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
	})
}

func (tcg *TaskConsumerGroup) getNewKafkaWrite(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(tcg.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(tcg.log.Debugf),
		ErrorLogger:  kafka.LoggerFunc(tcg.log.Errorf),
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
}

func (tcg *TaskConsumerGroup) consumeCreateTask(
	ctx context.Context,
	cancel context.CancelFunc,
	groupID string,
	topic string,
	workersNUm int,
) {
	r := tcg.getNewKafkaReader(tcg.Brokers, topic, groupID)
	defer cancel()
	defer func() {
		if err := r.Close(); err != nil {
			tcg.log.Error("r.Close", err)
			cancel()
		}
	}()

	w := tcg.getNewKafkaWrite(deadLetterQueueTopic)
	defer func() {
		if err := w.Close(); err != nil {
			tcg.log.Error("w.Close", err)
			cancel()
		}
	}()

	tcg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNUm; i++ {
		wg.Add(1)
		go tcg.createTaskWorker(ctx, cancel, r, w, wg, i)
	}

	wg.Wait()
}

func (tcg *TaskConsumerGroup) publishErrorMesssage(ctx context.Context, w *kafka.Writer, m kafka.Message, err error) error {
	errMsg := &model.ErrorMessage{
		Offset:    m.Offset,
		Error:     err.Error(),
		Time:      m.Time.UTC(),
		Partition: m.Partition,
		Topic:     m.Topic,
	}

	errMsgBytes, err := json.Marshal(errMsg)
	if err != nil {
		return err
	}

	return w.WriteMessages(ctx, kafka.Message{
		Value: errMsgBytes,
	})
}

func (tcg *TaskConsumerGroup) RunConsumers(ctx context.Context, cancel context.CancelFunc) {
	go tcg.consumeCreateTask(ctx, cancel, tasksGroupID, doneTaskTopic, taskWorkers)
}
