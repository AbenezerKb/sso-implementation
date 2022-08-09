Feature: Login

  Background:
    Given I am a registered user with details
      | phone         | email             | password |
      | 251911121314 | example@email.com | 1234abcd |

  @success
  Scenario Outline: Successful Login
    Given I fill the following details
      | phone   | email   | password   | otp   |
      | <phone> | <email> | <password> | <otp> |
    When I submit the registration form
    Then I will be logged in securely to my account
    Examples:
      | phone        | email             | password | otp  |
      | 251911121314 |                   |          | 123456 |
      |              | example@email.com | 1234abcd |      |
      | 251911121314 | example@email.com | 1234abcd | 123456 |

  @invalid
  Scenario Outline: Failed Login
    Given I fill the following details
      | phone   | email   | password   | otp   |
      | <phone> | <email> | <password> | <otp> |
    When I submit the registration form
    Then the login should fail with "<message>"
    Examples:
      | phone         | email             | password | otp    | message             |
      | +251911121314 |                   |          | 654321 | invalid credentials |
      |               | example@gmail.com | abcd1234 |        | invalid credentials |
      | +251914131211 |                   |          | 123456 | invalid credentials |
      |               | not@email.com     | 1234abcd |        | invalid credentials |