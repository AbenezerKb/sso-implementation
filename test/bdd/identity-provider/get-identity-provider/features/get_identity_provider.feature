Feature: Get Identity Provider

    As an admin,
    I want to get an identity provider
    So that I can view and manage the identity provider

    Background:
        Given I am logged in as admin user
            | email           | password      | role       |
            | admin@gmail.com | adminPassword | super-user |

    @success
    Scenario: Successful get identity provider
        Given There is identity provider with the following details
            | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |
        When I Get the identity provider
        Then I should successfully get the identity provider

    @failure
    Scenario Outline: Failed get identity provider
        Given I have identity provider with id "<id>"
        When I Get the identity provider
        Then Then I should get error with message "<message>"

        Examples:
            | id                                   | message                      |
            | 06f7404f-5402-4832-8b0b-53da2cdd7efc | no identity provider found   |
            | 3kjf0-kljkla0-afl30-afl-dk           | invalid identity provider id |