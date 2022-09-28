package process_events

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/test"
	"sync"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/segmentio/kafka-go"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type processMiniRideEventsTest struct {
	test.TestInstance
	apiTest        src.ApiTest
	StreamedEvents []kafkaEvent
	Users          []dto.User
}

type kafkaEvent struct {
	Event  string
	Driver request_models.Driver
}

func TestProcessMiniRideEvents(t *testing.T) {
	p := &processMiniRideEventsTest{}
	p.TestInstance = test.Initiate("../../../../")
	p.apiTest = src.ApiTest{}
	p.apiTest.InitializeTest(t, "process miniRide event's", "features/process_mini_ride_events.feature", p.InitializeScenario)
}

func (p *processMiniRideEventsTest) theyAreTheFollowingUsersOnSso(users *godog.Table) error {
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

func (p *processMiniRideEventsTest) miniRideStreamedTheFollowingEvents(events *godog.Table) error {
	eventsString, err := p.apiTest.ReadRows(events, []src.Type{
		{
			WithName: "driver",
			Kind:     src.Object,
			Columns:  []string{"first_name", "middle_name", "last_name", "phone", "profile_picture", "status", "swap_phones", "driverId", "id"},
		},
		{
			Column: "swap_phones",
			Kind:   src.Array,
		},
	}, true)
	if err != nil {
		return err
	}

	err = p.apiTest.UnmarshalJSON([]byte(eventsString), &p.StreamedEvents)
	if err != nil {
		return err
	}
	defer p.KafkaConn.Close()
	for i := 0; i < len(p.StreamedEvents); i++ {
		msg, err := json.Marshal(p.StreamedEvents[i].Driver)
		if err != nil {
			return err
		}

		_, err = p.KafkaConn.WriteMessages(kafka.Message{
			Key:   []byte(p.StreamedEvents[i].Event),
			Value: msg,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *processMiniRideEventsTest) iProcessThoseEvents() error {

	t := time.NewTicker(1 * time.Second)
	wg := new(sync.WaitGroup)

	for range t.C {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*2))

		msg, err := p.PlatformLayer.Kafka.ReadMessage(ctx)
		if err != nil {
			fmt.Println("error in kafka read message")
			break
		}
		wg.Add(1)
		go p.Module.Mini_rideModule.ProcessEvents(ctx, msg, wg)
	}

	wg.Wait()
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
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for i := 0; i < len(p.Users); i++ {
			_, _ = p.DB.DeleteUser(ctx, p.Users[i].ID)
		}
		for i := 0; i < len(p.StreamedEvents); i++ {
			if p.StreamedEvents[i].Event == "CREATE" {
				_, _ = p.DB.DeleteUser(ctx, p.StreamedEvents[i].Driver.ID)
			}
		}
		return ctx, nil
	})
	ctx.Step(`^I process those event\'s$`, p.iProcessThoseEvents)
	ctx.Step(`^mini ride streamed the following event\'s$`, p.miniRideStreamedTheFollowingEvents)
	ctx.Step(`^they are the following user\'s on sso$`, p.theyAreTheFollowingUsersOnSso)
	ctx.Step(`^they will have effect on following sso user\'s$`, p.theyWillHaveEffectOnFollowingSsoUsers)
}
