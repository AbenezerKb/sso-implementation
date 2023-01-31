Feature: update User Status

  As a admin
  I want to reset user's password
  So that the user gets a new password via text

  Background: I am logged in as admin
    Given I am logged in as admin user
      | email           | password | role                |
      | admin@gmail.com | iAmAdmin | reset_user_password |

  @success
  Scenario Outline: Successful User Status Update
    Given there is user with the following details:
      | first_name   | middle_name   | last_name   | phone   | email   | password   |
      | <first_name> | <middle_name> | <last_name> | <phone> | <email> | <password> |
    When I reset the user's password"
    Then the user's password should be changed

    Examples:
      | first_name | middle_name | last_name | phone        | email           | password |
      | testuser1  | testuser1   | testuser1 | 251925252525 | test1@gmail.com | 123456   |
