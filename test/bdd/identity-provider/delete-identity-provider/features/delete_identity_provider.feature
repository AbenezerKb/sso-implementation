Feature: Delete Identity Provider

    As an admin,
    I want to be able to delete identity providers
    So that I can remove unused and invalid providers

    Background:
        Given I am logged in as admin user
            | email           | password      | role                     |
            | admin@gmail.com | adminPassword | delete_identity_provider |

    @success
    Scenario: Successful delete identity provider
        Given There is identity provider with the following details
            | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |
        When I delete the identity provider
        Then The identity provider should be deleted

    @failure
    Scenario Outline: Failed delete identity provider
        Given There is identity provider with id "<id>"
        When I delete the identity provider
        Then The delete should fail with error message "<message>"

        Examples:
            | id                                 | message                     |
            | 60d56419-c2e9-4ee4-951f-04644d245e | identity provider not found |
            | 4                                  | identity provider not found |
