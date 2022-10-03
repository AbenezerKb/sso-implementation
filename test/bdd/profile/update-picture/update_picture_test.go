package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

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
	b, w := u.openMultipartFormData(picPath)
	u.apiTest.Body = b.String()
	u.apiTest.SetHeader("Content-Type", w.FormDataContentType())
	return nil
}

func (u *updateProfilePictureTest) iUpdateMyProfilePicture() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *updateProfilePictureTest) myProfilePictureShouldBeUpdated() error {
	updatedUser, err := u.DB.GetUserById(context.Background(), u.User.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(updatedUser.ProfilePicture.String, u.User.ProfilePicture.String); err == nil {
		return fmt.Errorf("profile picture not updated")
	}

	return nil
}

func (u *updateProfilePictureTest) theUpdateShouldFailWithMessage(arg1 string) error {
	return nil
}

func (u *updateProfilePictureTest) openMultipartFormData(filePath string) (*bytes.Buffer, *multipart.Writer) {

	file, _ := os.Open(filePath)
	defer file.Close()

	u.fileName = file.Name()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("image", file.Name())
	io.Copy(part, file)
	w.Close()

	return body, w
}

func (u *updateProfilePictureTest) InitializeScenario(ctx *godog.ScenarioContext) {
	u.apiTest.URL = "/v1/profile/picture"
	u.apiTest.Method = http.MethodPut
	u.apiTest.InitializeServer(u.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteUser(ctx, u.User.ID)

		return ctx, nil
	})

	ctx.Step(`^I am logged in user with the following details:$`, u.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I selected this picture "([^"]*)"$`, u.iSelectedThisPicture)
	ctx.Step(`^I update my profile picture$`, u.iUpdateMyProfilePicture)
	ctx.Step(`^my profile picture should be updated$`, u.myProfilePictureShouldBeUpdated)
	ctx.Step(`^The update should fail with message "([^"]*)"$`, u.theUpdateShouldFailWithMessage)
}
