package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateProfilePictureTest struct {
	test.TestInstance
	User        db.User
	NewUserData dto.User
	apiTest     src.ApiTest
	fileName    string
}

func TestUpdateProfilePicture(t *testing.T) {
	u := updateProfilePictureTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest.InitializeTest(t, "update profile picture", "features/update_picture.feature", u.InitializeScenario)
}

func (u *updateProfilePictureTest) iAmLoggedInUserWithTheFollowingDetails(userDetails *godog.Table) error {
	userData, err := u.apiTest.ReadRow(userDetails, nil, false)
	if err != nil {
		return err
	}

	userValue := dto.User{}
	err = u.apiTest.UnmarshalJSON([]byte(userData), &userValue)
	if err != nil {
		return err
	}

	u.User, err = u.AuthenticateWithParam(userValue)
	if err != nil {
		return err
	}
	u.apiTest.SetHeader("Authorization", "Bearer "+u.AccessToken)
	return nil
}

func (u *updateProfilePictureTest) iSelectedThisPicture(picPath string) error {
	u.apiTest.URL = "/v1/assets"
	u.apiTest.Method = http.MethodPost
	b, w := u.openMultipartFormData(picPath)
	u.apiTest.Body = b.String()
	u.apiTest.SetHeader("Content-Type", w.FormDataContentType())
	u.apiTest.SendRequest()

	if err := u.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return nil
	}

	imageName := struct {
		Data string `json:"data"`
	}{}
	err := u.apiTest.UnmarshalResponseBody(&imageName)
	if err != nil {
		return err
	}
	u.User.ProfilePicture.String = imageName.Data
	u.apiTest.ResetResponse()
	u.apiTest.SetBodyMap(map[string]interface{}{
		"first_name":      u.User.FirstName,
		"middle_name":     u.User.MiddleName,
		"last_name":       u.User.LastName,
		"gender":          u.User.Gender,
		"profile_picture": imageName.Data,
	})
	u.apiTest.SetHeader("Content-Type", "application/json")
	u.apiTest.URL = "/v1/profile"
	u.apiTest.Method = http.MethodPut
	u.apiTest.SendRequest()
	return nil
}

func (u *updateProfilePictureTest) iUpdateMyProfilePicture() error {
	// 	u.apiTest.SendRequest()
	return nil
}

func (u *updateProfilePictureTest) myProfilePictureShouldBeUpdated() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	updatedUser, err := u.DB.GetUserById(context.Background(), u.User.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(updatedUser.ProfilePicture.String, u.User.ProfilePicture.String); err != nil {
		return fmt.Errorf("profile picture not updated")
	}
	u.User = updatedUser

	return nil
}

func (u *updateProfilePictureTest) theUpdateShouldFailWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (u *updateProfilePictureTest) openMultipartFormData(filePath string) (*bytes.Buffer, *multipart.Writer) {

	file, _ := os.Open(filePath)
	defer file.Close()

	u.fileName = file.Name()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("asset", file.Name())
	io.Copy(part, file)
	_ = w.WriteField("type", "profile_picture")
	w.Close()

	return body, w
}

func (u *updateProfilePictureTest) InitializeScenario(ctx *godog.ScenarioContext) {
	u.apiTest.URL = "/v1/profile/picture"
	u.apiTest.Method = http.MethodPut
	u.apiTest.InitializeServer(u.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_ = os.Remove("../../../../static/profile_picture/" + u.User.ProfilePicture.String)

		_, _ = u.DB.DeleteUser(ctx, u.User.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in user with the following details:$`, u.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I selected this picture "([^"]*)"$`, u.iSelectedThisPicture)
	ctx.Step(`^I update my profile picture$`, u.iUpdateMyProfilePicture)
	ctx.Step(`^my profile picture should be updated$`, u.myProfilePictureShouldBeUpdated)
	ctx.Step(`^The update should fail with message "([^"]*)"$`, u.theUpdateShouldFailWithMessage)
}
