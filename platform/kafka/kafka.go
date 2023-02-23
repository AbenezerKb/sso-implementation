package kafkaconsumer

import (
	"context"
	"encoding/json"
	"log"

	"sso/internal/constant/errors"
	"sso/platform/logger"

	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(ctx context.Context, event json.RawMessage) error
type Kafka interface {
	RegisterKafkaEventHandler(EventType string, handler EventHandler)
	Close() error
}
type kafkaClient struct {
	kafkaReader   *kafka.Reader
	log           logger.Logger
	eventHandlers map[string]EventHandler
}

func NewKafkaConnection(kafkaURL, topic, groupID string, maxBytes int, logger logger.Logger) Kafka {
	_, err := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)
	if err != nil {
		logger.Fatal(context.Background(), "failed to connect kafka leader %v", zap.Error(err))
	}
	log.Printf("url: %v topic:  %v", kafkaURL, topic)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaURL},
		GroupID:   groupID,
		Topic:     topic,
		Partition: 0,
		MaxBytes:  maxBytes,
	})

	kafkaClient := &kafkaClient{
		log:           logger,
		kafkaReader:   r,
		eventHandlers: make(map[string]EventHandler),
	}

	go kafkaClient.readMessage(context.Background())
	return kafkaClient
}
func (k *kafkaClient) RegisterKafkaEventHandler(EventType string, handler EventHandler) {
	// event handlers for kafka event rypes
	h, ok := k.eventHandlers[EventType]
	if ok {
		k.log.Warn(context.Background(), "kafka event handler is being over written by a new event handler",
			zap.Any("old-handler:", h), zap.Any("new-handler:", handler))
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
		k.Close()
	}()

	// Loop Forever
	for {
		log.Print("reading message....\n")
		payload, err := k.kafkaReader.FetchMessage(ctx)
		if err != nil {
			k.log.Info(ctx, "kafka connection error", zap.Error(err), zap.Error(err))
			return
		}
		if payload.Value == nil {
			k.log.Warn(ctx, "kafka sent empty message", zap.Any("key:", payload.Key))
			continue
		}
		log.Printf("kafka event key %v : message %v", string(payload.Key), string(payload.Value))
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
