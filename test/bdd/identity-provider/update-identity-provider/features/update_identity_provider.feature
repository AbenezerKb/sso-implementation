Feature: Update Identity Provider

    As an admin,
    I want to update identity providers
    So that I can manage changes to URLs and other changes made by the providers

    Background:
        Given I am logged in as admin user
            | email           | password      | role       |
            | admin@gmail.com | adminPassword | super-user |
        And There is identity provider with the following details
            | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |


    @success
    Scenario Outline: Successful idP Update
        Given I fill the form with the following details
            | name   | logo_uri   | client_id   | client_secret   | redirect_uri   | authorization_uri   | token_endpoint_uri   | user_info_endpoint_uri   |
            | <name> | <logo_uri> | <client_id> | <client_secret> | <redirect_uri> | <authorization_uri> | <token_endpoint_uri> | <user_info_endpoint_uri> |
        When I update the identity provider
        Then The identity provider should be updated

        Examples:
            | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | http://authorize.com | http://token.com   | http.//userinfo.com    |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | http://authorize.com | http://token.com   |                        |

    @failure
    Scenario Outline: Failed idP Update
        Given I fill the form with the following details
            | name   | logo_uri   | client_id   | client_secret   | redirect_uri   | authorization_uri   | token_endpoint_uri   | user_info_endpoint_uri   |
            | <name> | <logo_uri> | <client_id> | <client_secret> | <redirect_uri> | <authorization_uri> | <token_endpoint_uri> | <user_info_endpoint_uri> |
        When I update the identity provider
        Then The identity provider updated should fail with message "<message>"

        Examples:
            | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri | message                        |
            |      | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | name is required               |
            | ip_1 | https://www.google.com/images/errors/robot.png |           | secret_2      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | client_id is required          |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  |               | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | client_secret is required      |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      |                      | http://authorize.com | http://token.com   | http.userinfo.com      | redirect_uri is required       |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | not-valid-uri        | http://authorize.com | http://token.com   | http.userinfo.com      | invalid redirect_uri           |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com |                      | http://token.com   | http.userinfo.com      | authorization_uri is required  |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | not-valid-uri        | http://token.com   | http.userinfo.com      | invalid authorization_uri      |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | http://authorize.com |                    | http.userinfo.com      | token_endpoint_uri is required |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | http://authorize.com | not-valid-uri      | http.userinfo.com      | invalid token_endpoint_uri     |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | http://authorize.com | http://token.com   | not-valid-uri          | invalid user_info_endpoint_uri |
            | ip_1 | not-valid-uri                                  | client_2  | secret_2      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | invalid logo_uri               |
            | ip_1 | https://www.not-found-logo.com/logo.png        | client_2  | secret_2      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | logo not found                 |
