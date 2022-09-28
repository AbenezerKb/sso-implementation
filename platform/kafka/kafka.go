package kafka_consumer

import (
	"context"
	"encoding/json"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto/request_models"
	"sso/platform"
	"sso/platform/logger"
	"strings"

	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type kafkaClient struct {
	kafkaURL    string
	topic       string
	kafkaReader *kafka.Reader
	log         logger.Logger
	groupID     string
}

func NewKafkaConnection(kafkaURL, topic, groupID string, log logger.Logger) platform.Kafka {
	kafkaClient := kafkaClient{
		kafkaURL: kafkaURL,
		topic:    topic,
		groupID:  groupID,
		log:      log,
	}
	kafkaReader := kafkaClient.getKafkaReader()
	kafkaClient.kafkaReader = kafkaReader
	return &kafkaClient
}

func (k *kafkaClient) getKafkaReader() *kafka.Reader {
	brokers := strings.Split(k.kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   k.topic,
		GroupID: k.groupID,
	})
}

func (k *kafkaClient) Close() error {
	return k.kafkaReader.Close()
}

func (k *kafkaClient) ReadMessage(ctx context.Context) (*request_models.MinRideEvent, error) {
	var rsp request_models.MinRideEvent
	msg, err := k.kafkaReader.ReadMessage(ctx)
	if err != nil {
		err = errors.ErrKafkaRead.Wrap(err, "couldn't read kafka event")
		k.log.Debug(context.Background(), "couldn't read message from kafka", zap.Error(err))
		return nil, err
	}

	rsp.Event = string(msg.Key)

	err = json.Unmarshal(msg.Value, &rsp.Driver)
	if err != nil {
		err = errors.ErrKafkaInvalidEvent.Wrap(err, "couldn't unmarshal kafka key")
		k.log.Error(context.Background(), "couldn't unmarshal kafka value", zap.Any("value", msg.Value))
		return nil, err
	}
	k.log.Info(ctx, "successFully read message from kafka", zap.Any("msg", rsp))

	return &rsp, nil
}
