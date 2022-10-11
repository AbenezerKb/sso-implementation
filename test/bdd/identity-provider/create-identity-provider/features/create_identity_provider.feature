Feature: Create Identity Provider
  As an admin,
  I want to create an identity provider
  So that users can login with other identity providers

  Background:
    Given I am logged in with the following credentials
      | email           | password | role                     |
      | admin@gmail.com | 12345678 | create_identity_provider |

  @success
  Scenario Outline: I create identity provider successfully
    Given I have filled the following data for the identity provider
      | name   | logo_uri   | client_id   | client_secret   | redirect_uri   | authorization_uri   | token_endpoint_uri   | user_info_endpoint_uri   |
      | <name> | <logo_uri> | <client_id> | <client_secret> | <redirect_uri> | <authorization_uri> | <token_endpoint_uri> | <user_info_endpoint_uri> |
    When I submit to create an identity provider
    Then the identity provider should be created
    Examples:
      | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   |                        |

  @failure
  Scenario Outline: I fail to create identity provider
    Given I have filled the following data for the identity provider
      | name   | logo_uri   | client_id   | client_secret   | redirect_uri   | authorization_uri   | token_endpoint_uri   | user_info_endpoint_uri   |
      | <name> | <logo_uri> | <client_id> | <client_secret> | <redirect_uri> | <authorization_uri> | <token_endpoint_uri> | <user_info_endpoint_uri> |
    When I submit to create an identity provider
    Then the request should fail with "<message>"
    Examples:

      | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri | message                        |
      |      | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | name is required               |
      | ip_1 | https://www.google.com/images/errors/robot.png |           | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | client_id is required          |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  |               | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | client_secret is required      |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      |                      | http://authorize.com | http://token.com   | http.userinfo.com      | redirect_uri is required       |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | not-valid-uri        | http://authorize.com | http://token.com   | http.userinfo.com      | invalid redirect_uri           |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com |                      | http://token.com   | http.userinfo.com      | authorization_uri is required     |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | not-valid-uri        | http://token.com   | http.userinfo.com      | invalid authorization_uri      |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com |                    | http.userinfo.com      | token_endpoint_uri is required |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | not-valid-uri      | http.userinfo.com      | invalid token_endpoint_uri     |
      | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | not-valid-uri          | invalid user_info_endpoint_uri |
      | ip_1 | not-valid-uri                                  | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | invalid logo_uri               |
      | ip_1 | https://www.not-found-logo.com/logo.png        | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      | logo not found                 |
