package process_events

import (
	"context"
	"database/sql"
	"encoding/json"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/test"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/segmentio/kafka-go"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type processMiniRideEventsTest struct {
	test.TestInstance
	apiTest src.ApiTest
	Users   []dto.User
}

func TestProcessMiniRideEvents(t *testing.T) {
	p := &processMiniRideEventsTest{}
	p.TestInstance = test.Initiate("../../../../")
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
	err = p.KafkaWriter.WriteMessages(context.Background(), messages...)
	if err != nil {

		return err
	}
	return nil
}

func (p *processMiniRideEventsTest) iProcessThoseEvents() error {
	time.Sleep(time.Second)
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
