package kafkaconsumer

import (
	"context"
	"encoding/json"
	"time"

	"sso/internal/constant/errors"
	"sso/platform/logger"
	"sso/platform/routine"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type EventHandler func(ctx context.Context, event json.RawMessage) error
type Kafka interface {
	RegisterKafkaEventHandler(EventType string, handler EventHandler)
	Close() error
}
type ContextKey any

type kafkaClient struct {
	kafkaReader   *kafka.Reader
	log           logger.Logger
	maxBytes      int
	eventHandlers map[string]EventHandler
}

func NewKafkaConnection(kafkaURL, groupID string, topics []string, maxBytes int,
	log logger.Logger) Kafka {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaURL},
		GroupID:     groupID,
		GroupTopics: topics,
		MinBytes:    10,
		MaxBytes:    maxBytes,
		ErrorLogger: log.Named("kafka-reader"),
	})
	kafkaClient := kafkaClient{
		kafkaReader:   r,
		log:           log,
		maxBytes:      maxBytes,
		eventHandlers: make(map[string]EventHandler),
	}
	// run the read message
	routine.ExecuteRoutine(context.Background(), routine.Routine{
		Name: "kafka-read-message",
		Operation: func(ctx context.Context, log logger.Logger) {
			kafkaClient.readMessage(context.Background())
		},
		NoWait: true,
	}, log)
	return &kafkaClient
}
func (k *kafkaClient) RegisterKafkaEventHandler(EventType string, handler EventHandler) {
	// register event handlers for kafka event types
	h, ok := k.eventHandlers[EventType]
	if ok {
		k.log.Warn(context.Background(), "kafka event handler is being over written by a new event handler",
			zap.Any("event-type", EventType), zap.Any("old-handler:", h), zap.Any("new-handler:", handler))
	}
	k.eventHandlers[EventType] = handler
}

func (k *kafkaClient) Close() error {
	return k.kafkaReader.Close()
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (k *kafkaClient) routeEvent(ctx context.Context, payload kafka.Message) error {
	// Marshal incoming data into a Event struct
	var request json.RawMessage
	if err := json.Unmarshal(payload.Value, &request); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid message data")
		k.log.Info(ctx, "tried to read invalid message data", zap.Error(err))
		return err
	}
	handler, ok := k.eventHandlers[string(payload.Key)]
	if !ok {
		err := errors.ErrKafkaInvalidEvent.New("this event does not have a handler set:%v", string(payload.Key))
		k.log.Warn(ctx, "unsupported event by kafka handlers", zap.Error(err))
		return err
	}

	ctx = context.WithValue(
		context.WithValue(
			context.Background(),
			ContextKey("request-start-time"), time.Now()),
		ContextKey("x-request-id"), uuid.NewString())
	// Execute the handler and return any nil
	err := handler(ctx, request)
	if err != nil {
		return err
	}
	return nil

}

// this should be run on a go routine
func (k *kafkaClient) readMessage(ctx context.Context) {
	// read the message and route to the correct handler
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		err := k.Close()
		if err != nil {
			k.log.Warn(ctx, "error while closing kafka connection")
		}
	}()

	// Loop Forever
	for {
		payload, err := k.kafkaReader.FetchMessage(ctx)
		if err != nil {
			k.log.Info(ctx, "kafka connection error", zap.Error(err), zap.Error(err))
			return
		}
		if payload.Value == nil {
			k.log.Warn(ctx, "kafka sent empty message", zap.Any("key:", payload.Key))
			continue
		}
		if err := k.routeEvent(ctx, payload); err != nil {
			k.log.Warn(ctx, "event handler faild to process kafka request", zap.Error(err))
		} else {
			err := k.kafkaReader.CommitMessages(ctx, payload)
			if err != nil {
				k.log.Warn(ctx, "failed to commit processed message", zap.Error(err))
			}
		}

	}
}
