package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"

	"youtube-downloader-rest/config"
	"youtube-downloader-rest/pkg/logger"
)

type DoneTaskProducer interface {
	PublishCreate(ctx context.Context, msgs ...kafka.Message) error
	PublishUpdate(ctx context.Context, msgs ...kafka.Message) error
	Close()
	Run()
	GetNewKafkaWriter(topic string) *kafka.Writer
}

type doneTaskProducer struct {
	log          logger.Logger
	cfg          *config.Config
	createWriter *kafka.Writer
	updateWriter *kafka.Writer
}

func NewDoneTaskProducer(log logger.Logger, cfg *config.Config) *doneTaskProducer {
	return &doneTaskProducer{log: log, cfg: cfg}
}

func (t *doneTaskProducer) GetNewKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(t.cfg.Kafka.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(t.log.Debugf),
		ErrorLogger:  kafka.LoggerFunc(t.log.Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
}

func (t *doneTaskProducer) Run() {
	t.createWriter = t.GetNewKafkaWriter(doneTaskTopic)
}

func (t *doneTaskProducer) Close() {
	t.createWriter.Close()
	t.updateWriter.Close()
}

func (t *doneTaskProducer) PublishCreate(ctx context.Context, msgs ...kafka.Message) error {
	return t.createWriter.WriteMessages(ctx, msgs...)
}

func (t *doneTaskProducer) PublishUpdate(ctx context.Context, msgs ...kafka.Message) error {
	return t.updateWriter.WriteMessages(ctx, msgs...)
}
