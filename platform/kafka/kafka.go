package kafkaconsumer

import (
	"context"
	"encoding/json"

	"sso/internal/constant/errors"
	"sso/internal/storage"
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
	kafkaURL      string
	topic         string
	kafkaConn     *kafka.Conn
	log           logger.Logger
	offsetStore   storage.Kafka
	groupID       string
	maxBytes      int
	eventHandlers map[string]EventHandler
}

func NewKafkaConnection(kafkaURL, topic, groupID string, maxBytes int, log logger.Logger, offsetPersitence storage.Kafka) Kafka {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)
	if err != nil {
		log.Error(context.Background(), "failed not dail kafka leader", zap.Error(err))
		return nil
	}

	kafkaClient := &kafkaClient{
		kafkaURL:      kafkaURL,
		topic:         topic,
		groupID:       groupID,
		log:           log,
		kafkaConn:     conn,
		offsetStore:   offsetPersitence,
		maxBytes:      maxBytes,
		eventHandlers: make(map[string]EventHandler),
	}
	offset, err := kafkaClient.offsetStore.GetOffset(context.Background())
	if err != nil {
		log.Fatal(context.Background(), "unable to get offset for kafka connection", zap.Error(err))
	}
	_, err = kafkaClient.kafkaConn.Seek(offset, kafka.SeekStart)
	if err != nil {
		log.Fatal(context.Background(), "unable to get offset for kafka connection", zap.Error(err))
	}
	// run the read message
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
	return k.kafkaConn.Close()
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
		payload, err := k.kafkaConn.ReadMessage(k.maxBytes)
		if err != nil {
			k.log.Info(ctx, "kafka connection error", zap.Error(err), zap.Error(err))
			return
		}
		if payload.Value == nil {
			k.log.Warn(ctx, "kafka sent empty message", zap.Any("key:", payload.Key))
			continue
		}
		err = k.offsetStore.SetOffset(ctx, payload.Offset+1)
		if err != nil {
			err = errors.ErrInternalServerError.Wrap(err, "faild to set offset for sso consumer")
			k.log.Error(ctx, "faild to set consumer offset", zap.Error(err))
		}
		if err := k.routeEvent(ctx, payload); err != nil {
			k.log.Warn(ctx, "event handler faild to process kafka request", zap.Error(err))
		}

	}
}
