package mini_ride

import (
	"context"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"
	"sync"

	"github.com/dongri/phonenumber"
	"go.uber.org/zap"
)

type miniRide struct {
	log                 logger.Logger
	miniRidePersistence storage.MiniRidePersistence
	kafkaClient         platform.Kafka
}

func InitMinRide(log logger.Logger, miniRidePersistence storage.MiniRidePersistence, kafkaClient platform.Kafka) module.MiniRideModule {
	return &miniRide{
		log:                 log,
		miniRidePersistence: miniRidePersistence,
		kafkaClient:         kafkaClient,
	}

	// go m.listenMiniRideEvent(context.Background())
}

func (m *miniRide) ListenMiniRideEvent(ctx context.Context) {
	wg := new(sync.WaitGroup)
	for {
		miniRideEvent, err := m.kafkaClient.ReadMessage(ctx)
		if err != nil {
			m.log.Error(ctx, "error reading kafka message", zap.Error(err))
			break
		}
		wg.Add(1)
		go m.ProcessEvents(ctx, miniRideEvent, wg)
	}
	wg.Wait()
}

func (m *miniRide) ProcessEvents(ctx context.Context, miniRideEvent *request_models.MinRideEvent, wg *sync.WaitGroup) {
	defer wg.Done()
	switch miniRideEvent.Event {
	case constant.UPDATE:
		if len(miniRideEvent.Driver.SwapPhones) > 1 {
			// swap phone
			err := m.miniRidePersistence.SwapPhones(ctx, miniRideEvent.Driver.SwapPhones[0], miniRideEvent.Driver.SwapPhones[1])
			if err != nil {
				return
			}
		}
		// update user
		err := m.miniRidePersistence.UpdateUser(ctx, miniRideEvent.Driver)
		if err != nil {
			m.log.Error(ctx, "couldn't update user", zap.Any("user", miniRideEvent.Driver))
			return
		}
	case constant.CREATE:
		// create user
		_, err := m.miniRidePersistence.CreateUser(ctx, miniRideEvent.Driver)
		if err != nil {
			return
		}
	case constant.PROMOTE:
		// update user
		err := m.miniRidePersistence.UpdateUser(ctx, miniRideEvent.Driver)
		if err != nil {
			m.log.Error(ctx, "couldn't update user", zap.Any("user", miniRideEvent.Driver))
			return
		}

	default:
		m.log.Debug(ctx, "unwanted event form kafka", zap.Any("event", miniRideEvent.Event))
	}

}

func (m *miniRide) CheckPhone(ctx context.Context, phone string) (*dto.MiniRideResponse, error) {

	parsedPhone := phonenumber.Parse(phone, "ET")
	if parsedPhone == "" {
		err := errors.ErrInvalidUserInput.New("invalid phone number")
		m.log.Error(ctx, "couldn't parse phone", zap.Error(err), zap.String("phone", phone))
		return nil, err
	}

	return m.miniRidePersistence.CheckPhone(ctx, parsedPhone)
}
