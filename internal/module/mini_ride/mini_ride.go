package mini_ride

import (
	"context"
	"encoding/json"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/module"
	"sso/internal/storage"
	kafka_consumer "sso/platform/kafka"
	"sso/platform/logger"
	"strings"

	"github.com/dongri/phonenumber"
	"go.uber.org/zap"
)

type miniRide struct {
	log                 logger.Logger
	miniRidePersistence storage.MiniRidePersistence
	kafkaClient         kafka_consumer.Kafka
}

func InitMinRide(log logger.Logger, miniRidePersistence storage.MiniRidePersistence, kafkaClient kafka_consumer.Kafka) module.MiniRideModule {
	return &miniRide{
		log:                 log,
		miniRidePersistence: miniRidePersistence,
		kafkaClient:         kafkaClient,
	}

	// go m.listenMiniRideEvent(context.Background())
}

func (m *miniRide) parseRideMiniResponse(ctx context.Context, rideMiniData []byte) (*request_models.Driver, error) {
	// marshal the data to rideminiresponse
	rideMiniResponse := request_models.MiniRideDriverResponse{}
	err := json.Unmarshal(rideMiniData, &rideMiniResponse)
	if err != nil {
		m.log.Error(ctx, "unable to bind ridemini data to local dto", zap.Error(err))
		err = errors.ErrInvalidUserInput.Wrap(err, "unable to bind ridemini dataunable to bind ridemini data ")
		return nil, err
	}
	names := strings.Split(rideMiniResponse.FullName, " ")
	result := &request_models.Driver{
		ID:             rideMiniResponse.ID,
		DriverID:       rideMiniResponse.DriverID,
		Phone:          rideMiniResponse.Phone,
		Status:         rideMiniResponse.Status,
		ProfilePicture: rideMiniResponse.ProfilePicture,
		SwapPhones:     rideMiniResponse.SwapPhones,
	}
	for i := range names {
		if i == 0 {
			result.FirstName = names[0]
		}
		if i == 1 {
			result.MiddleName = names[1]
		}
		if i == 2 {
			result.LastName = names[2]
		}
	}
	return result, nil
}
func (m *miniRide) CreateUser(ctx context.Context, data json.RawMessage) error {
	driver, err := m.parseRideMiniResponse(ctx, data)
	if err != nil {
		return err
	}
	// create user
	_, err = m.miniRidePersistence.CreateUser(ctx, driver)
	if err != nil {
		return nil
	}
	return nil
}
func (m *miniRide) UpdateUser(ctx context.Context, data json.RawMessage) error {
	driver, err := m.parseRideMiniResponse(ctx, data)
	if err != nil {
		return err
	}
	if len(driver.SwapPhones) > 1 {
		// swap phone
		err := m.miniRidePersistence.SwapPhones(ctx, driver.SwapPhones[0], driver.SwapPhones[1])
		if err != nil {
			return err
		}
	}
	// update user
	err = m.miniRidePersistence.UpdateUser(ctx, driver)
	if err != nil {
		m.log.Error(ctx, "couldn't update user", zap.Any("user", driver))
		return err
	}
	return nil
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
