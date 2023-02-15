package kafka

import (
	"context"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"
)

type kafkaPersistence struct {
	logger logger.Logger
	db     *persistencedb.PersistenceDB
}

func InitKafkaPersistence(logger logger.Logger, db *persistencedb.PersistenceDB) storage.Kafka {
	return &kafkaPersistence{
		logger: logger,
		db:     db,
	}
}
func (k *kafkaPersistence) GetOffset(ctx context.Context) (int64, error) {
	offset, err := k.db.GetKafkaOffset(ctx)
	if err != nil {
		return 0, err
	}
	return int64(offset), nil
}
func (k *kafkaPersistence) SetOffset(ctx context.Context, offset int64) error {
	return k.db.SetKafkaOffset(ctx, int32(offset))
}
