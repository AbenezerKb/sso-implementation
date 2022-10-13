Feature: Login with Identity Provider
  As a user,
  I want to login with other SSO servers
  So that I can use my existing credentials to login to this SSO

  Background:
    Given There exists an identity provider with the following info
      | name | client_id | client_secret | token_endpoint_url      |
      | ip_1 | some_id   | some_secret   | https://token.com/token |
    And I am registered on that identity provider as follows
      | id    | first_name | middle_name | last_name | phone      | email         | gender | profile_picture              |
      | my-id | Trent      | Alexander   | Arnold    | 0912233445 | taa@gmail.com | male   | https://logo.com/picture.png |

  Scenario: I successfully login with the identity provider
    Given I have granted consent to my login with code "veryLegitCode"
    When I request to login with identity provider "ip_1"
    Then I should successfully login

  Scenario Outline: I fail to login because of bad request
    Given I have granted consent to my login with code "<code>"
    When I request to login with identity provider "<ip>"
    Then my request should fail with message "<message>"
    Examples:
      | code          | ip   | message                       |
      |               | ip_1 | code is required              |
      | veryLegitCode |      | identity provider is required |

  Scenario Outline: I fail to login because of invalid data
    Given I have granted consent to my login with code "<code>"
    When I request to login with identity provider "<ip>"
    Then my request should fail with "<message>"
    Examples:
      | code          | ip                                   | message                                                                       |
      | invalid-code  | ip_1                                 | authentication failed                                                         |
      | veryLegitCode | ip_2                                 | invalid identity provider                                                     |
      | veryLegitCode | 18760bc8-ffc2-405d-963e-e9ea6c3cd36b | identity provider with id 18760bc8-ffc2-405d-963e-e9ea6c3cd36b does not exist |
