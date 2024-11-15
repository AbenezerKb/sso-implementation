package process_events

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	kafkaconsumer "sso/platform/kafka"
	"sso/test"
	"strconv"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/segmentio/kafka-go"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type processMiniRideEventsTest struct {
	test.TestInstance
	apiTest   src.ApiTest
	Users     []dto.User
	AfterFunc func()
}

func TestProcessMiniRideEvents(t *testing.T) {
	p := &processMiniRideEventsTest{}

	p.TestInstance = test.Initiate("../../../../")

	conn, err := kafka.Dial("tcp", p.KafkaBroker)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("failed to get the current kafka controller:", err)
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatal("failed to dail kafka leader via non leader broker:", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             p.KafkaTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		if err != kafka.TopicAlreadyExists {
			log.Fatal("failed to create topics:", err)
		}
	}

	p.KafkaInitiator = kafkaconsumer.NewKafkaConnection(p.KafkaBroker, p.KafkaGroupID, []string{p.KafkaTopic}, p.KafkaMaxBytes, p.KafkaLogger)
	p.KafkaInitiator.RegisterKafkaEventHandler(string("CREATE"), p.Module.MiniRideModule.CreateUser)
	p.KafkaInitiator.RegisterKafkaEventHandler(string("UPDATE"), p.Module.MiniRideModule.UpdateUser)
	p.KafkaWritter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{p.KafkaBroker},
		Topic:        p.KafkaTopic,
		RequiredAcks: -1, //leading broker should acknowledge
		Logger:       p.KafkaLogger.Named("kafka-writter"),
		ErrorLogger:  p.KafkaLogger.Named("kafka-wrriter-errors"),
	})
	p.apiTest.InitializeServer(p.Server)
	p.apiTest.RunTest(t,
		"sync sso with ride-mini",
		&src.TestOptions{
			Paths: []string{"features/process_mini_ride_events.feature"},
		},
		p.InitializeScenario,
		func(ctx *godog.TestSuiteContext) {
			ctx.AfterSuite(func() {
				if err := p.DBCleanUp(); err != nil {
					t.Error(err)
				}
			})
		},
	)
}

func (p *processMiniRideEventsTest) thereAreTheFollowingUserDataOnSso(users *godog.Table) error {
	usersJson, err := p.apiTest.ReadRows(users, nil, false)
	if err != nil {
		return err
	}

	err = p.apiTest.UnmarshalJSON([]byte(usersJson), &p.Users)
	if err != nil {
		return err
	}

	for i := 0; i < len(p.Users); i++ {
		user, err := p.DB.CreateUserWithID(context.Background(), db.CreateUserWithIDParams{
			FirstName:      p.Users[i].FirstName,
			MiddleName:     p.Users[i].MiddleName,
			LastName:       p.Users[i].LastName,
			Phone:          p.Users[i].Phone,
			ProfilePicture: sql.NullString{String: p.Users[i].ProfilePicture, Valid: true},
			ID:             p.Users[i].ID,
		})
		if err != nil {
			return err
		}

		_, err = p.DB.UpdateUser(context.Background(), db.UpdateUserParams{
			Status: sql.NullString{String: p.Users[i].Status, Valid: true},
			ID:     user.ID,
		})
		if err != nil {
			return err
		}

		p.Users[i] = dto.User{
			ID:             user.ID,
			FirstName:      user.FirstName,
			MiddleName:     user.MiddleName,
			LastName:       user.LastName,
			Phone:          user.Phone,
			Status:         user.Status.String,
			ProfilePicture: user.ProfilePicture.String,
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *processMiniRideEventsTest) miniRideStreamedTheFollowingEvents(rideMiniData *godog.Table) error {
	defer p.KafkaWritter.Close()
	rows, err := p.apiTest.ReadRows(rideMiniData, []src.Type{
		{
			Column: "event",
			Ignore: true,
		},
		{
			Column: "swap_phones",
			Kind:   src.Array,
		},
		{
			Column: "id",
			Kind:   src.String,
		},
		{
			Column: "full_name",
			Kind:   src.String,
		},

		{
			Column: "driver_id",
			Kind:   src.String,
		},
		{
			Column: "driver_license",
			Kind:   src.String,
		},
		{
			Column: "phone",
			Kind:   src.String,
		},
		{
			Column: "profile_picture",
			Kind:   src.String,
		},
		{
			Column: "status",
			Kind:   src.String,
		},
	}, false)
	if err != nil {

		return err
	}
	eventsString, err := p.apiTest.ReadRows(rideMiniData, []src.Type{
		{
			Column: "event",
			Kind:   src.String,
		}}, false)
	if err != nil {

		return err
	}

	var eventStringArray []struct {
		Event string `json:"event"`
	}

	err = p.apiTest.UnmarshalJSON([]byte(eventsString), &eventStringArray)
	if err != nil {

		return err
	}
	rideMiniDrivers := []request_models.MiniRideDriverResponse{}
	err = p.apiTest.UnmarshalJSON([]byte(rows), &rideMiniDrivers)
	if err != nil {

		return err
	}
	messages := []kafka.Message{}
	for i := 0; i < len(rideMiniDrivers); i++ {
		rideminiDriver, err := json.Marshal(rideMiniDrivers[i])
		if err != nil {

			return err
		}
		messages = append(messages, kafka.Message{
			Key:   []byte(eventStringArray[i].Event),
			Value: rideminiDriver,
		})
	}
	err = p.KafkaWritter.WriteMessages(context.Background(), messages...)
	if err != nil {
		return err
	}
	return nil
}

func (p *processMiniRideEventsTest) iProcessThoseEvents() error {
	time.Sleep(10 * time.Second)
	return nil
}
func (p *processMiniRideEventsTest) theyWillHaveEffectOnFollowingSsoUsers(users *godog.Table) error {
	usersJson, err := p.apiTest.ReadRows(users, nil, false)
	if err != nil {
		return err
	}

	var usersStruct []dto.User
	err = p.apiTest.UnmarshalJSON([]byte(usersJson), &usersStruct)
	if err != nil {
		return err
	}
	for i := 0; i < len(usersStruct); i++ {
		fetchedUser, err := p.DB.GetUserById(context.Background(), usersStruct[i].ID)
		if err != nil {
			return err
		}

		if err := p.apiTest.AssertEqual(fetchedUser.FirstName, usersStruct[i].FirstName); err != nil {
			return err
		}
		if err := p.apiTest.AssertEqual(fetchedUser.MiddleName, usersStruct[i].MiddleName); err != nil {
			return err
		}
		if err := p.apiTest.AssertEqual(fetchedUser.LastName, usersStruct[i].LastName); err != nil {
			return err
		}
		if err := p.apiTest.AssertEqual(fetchedUser.ProfilePicture.String, usersStruct[i].ProfilePicture); err != nil {
			return err
		}
		if err := p.apiTest.AssertEqual(fetchedUser.Status.String, usersStruct[i].Status); err != nil {
			return err
		}
		if err := p.apiTest.AssertEqual(fetchedUser.Phone, usersStruct[i].Phone); err != nil {
			return err
		}
	}

	return nil
}

func (p *processMiniRideEventsTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I process those event\'s$`, p.iProcessThoseEvents)
	ctx.Step(`^mini ride streamed the following event\'s$`, p.miniRideStreamedTheFollowingEvents)
	ctx.Step(`^there are the following user data on sso$`, p.thereAreTheFollowingUserDataOnSso)
	ctx.Step(`^they will have effect on following sso user\'s$`, p.theyWillHaveEffectOnFollowingSsoUsers)

}
