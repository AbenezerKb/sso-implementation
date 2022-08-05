Feature: Registration

  Scenario: Successful Registration
    When I fill the form with the following details
      | first_name | middle_name | last_name | phone      | email           | password |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com | 1234567  |
    And I submit the registration form
    Then I will have a new account

  Scenario Outline: Failed Registration
    When I fill the form with the following details
      | first_name   | middle_name   | last_name   | phone   | email   | password   |
      | <first_name> | <middle_name> | <last_name> | <phone> | <email> | <password> |
    And I submit the registration form
    Then the registration should fail with "<message>"

    Examples:
      | first_name | middle_name | last_name | phone      | email           | password | message                                      |
      |            | testuser1   | testuser1 | 0925252525 | test1@gmail.com | 1234567  | first name is required                       |
      | testuser1  |             | testuser1 | 0925252525 | test1@gmail.com | 1234567  | middle name is required                      |
      | testuser1  | testuser1   |           | 0925252525 | test1@gmail.com | 1234567  | last name is required                        |
      | testuser1  | testuser1   | testuser1 |            | test1@gmail.com | 1234567  | phone is required                            |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com |          | password is required                         |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1gmail.com  | 1234567  | email is not valid                           |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com | 1jkl2    | password must be between 6 and 32 characters |
      | testuser1  | testuser1   | testuser1 | 33333333   | test1@gmail.com | 1234567  | invalid phone number                         |