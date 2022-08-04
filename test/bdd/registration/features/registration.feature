Feature: Registration
    Scenario: Successful Registration
        When I fill the form with the following details
            | first_name | middle_name | last_name | phone_number | email           | password |
            | testuser1  | testuser1   | testuser1 | 0925252525   | test1@gmail.com | 1234567  |
        And I submit the registration form
        Then I will have a new account

    Scenario Outline: Failed Registration
        When I fill the form with the following details
            | first_name   | middle_name   | last_name   | phone_number   | email   | password   |
            | <first_name> | <middle_name> | <last_name> | <phone_number> | <email> | <password> |
        And I submit the registration form
        Then the registration should fail with "<message>"

        Examples:
            | first_name | middle_name | last_name | phone_number | email           | password | message                  |
            |            | testuser1   | testuser1 | 0925252525   | test1.com       | 1234567  | First name is required   |
            | testuser1  |             | testuser1 | 0925252525   | test2@gmail.com | 1234567  | Middle name is required  |
            | testuser1  | testuser1   |           | 0925252525   | test1@gmail.com | 1234567  | Last name is required    |
            | testuser1  | testuser1   | testuser1 |              | test1@gmail.com | 1234567  | Phone number is required |
            | testuser1  | testuser1   | testuser1 | 0925252525   |                 | 1234567  | email is required        |
            | testuser1  | testuser1   | testuser1 | 0925252525   | test1@gmail.com |          | Password is required     |
            | testuser1  | testuser1   | testuser1 | 0925252525   | test1gmail.com  | 1234567  | Invalid email            |
            | testuser1  | testuser1   | testuser1 | 0925252525   | test1@gmail.com | 12       | Password to short        |
            | testuser1  | testuser1   | testuser1 | 52525        | test1@gmail.com | 1234567  | Invalid phone number     |