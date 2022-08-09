Feature: Registration

  Scenario: Successful Registration
    When I fill the form with the following details
      | first_name | middle_name | last_name | phone      | email           | password | otp    |
      | testuser1  | testuser1   | testuser1 | 0925252595 | test11@gmail.com | 1234567  | 123456 |
    And I submit the registration form
    Then I will have a new account

  Scenario Outline: Failed Registration
    When I fill the form with the following details
      | first_name   | middle_name   | last_name   | phone   | email   | password   | otp   |
      | <first_name> | <middle_name> | <last_name> | <phone> | <email> | <password> | <otp> |
    And I submit the registration form
    Then the registration should fail with "<message>"

    Examples:
      | first_name | middle_name | last_name | phone      | email           | password | otp    | message                                      |
      |            | testuser1   | testuser1 | 0925252525 | test1@gmail.com | 1234567  | 123456 | first name is required                       |
      | testuser1  |             | testuser1 | 0925252525 | test1@gmail.com | 1234567  | 123456 | middle name is required                      |
      | testuser1  | testuser1   |           | 0925252525 | test1@gmail.com | 1234567  | 123456 | last name is required                        |
      | testuser1  | testuser1   | testuser1 |            | test1@gmail.com | 1234567  | 123456 | phone is required                            |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com |          | 123456 | password is required                         |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1gmail.com  | 1234567  | 123456 | email is not valid                           |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com | 1jkl2    | 123456 | password must be between 6 and 32 characters |
      | testuser1  | testuser1   | testuser1 | 33333333   | test1@gmail.com | 1234567  | 123456 | invalid phone number                         |
      | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com | 1234567  | 12     | otp must be 6 characters                     |