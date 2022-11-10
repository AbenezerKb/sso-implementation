Feature: Client Registration

  Background: I am logged in as admin
    Given I am logged in as admin user
      | email           | password | role          |
      | admin@gmail.com | iAmAdmin | create_client |

  @success
  Scenario: Client Registers Successfully
    Given I fill the following client form
      | name      | client_type  | redirect_uris      | scopes        | logo_url                                       |
      | newClient | confidential | https://google.com | profile email | https://www.google.com/images/errors/robot.png |
    When I submit the form
    Then The registration should be successful

  @failure
  Scenario Outline: Client Registration Failure
    Given I fill the following client form
      | name   | client_type   | redirect_uris   | scopes   | logo_url   |
      | <name> | <client_type> | <redirect_uris> | <scopes> | <logo_url> |
    When I submit the form
    Then The registration should fail with "<message>"
    Examples:
      | name      | client_type  | redirect_uris           | scopes        | logo_url                                       | message                                           |
      |           | confidential | https://google.com      | profile email | https://www.google.com/images/errors/robot.png | name is required                                  |
      | newClient |              | https://google.com      | profile email | https://www.google.com/images/errors/robot.png | client_type is required                           |
      | newClient | confidential |                         | profile email | https://www.google.com/images/errors/robot.png | redirect_uris is required                         |
      | newClient | confidential | https://google.com      |               | https://www.google.com/images/errors/robot.png | scopes is required                                |
      | newClient | confidential | https://google.com      | profile email |                                                | logo_url is required                              |
      | newClient | my_type      | https://google.com      | profile email | https://www.google.com/images/errors/robot.png | client type must be either confidential or public |
#      | newClient | confidential | https://google.com      | not a scope   | https://www.google.com/images/errors/robot.png | invalid scopes                                    |
      | newClient | confidential | https://google.com      | profile email | my-logo-url                                    | invalid logo_url                                  |
