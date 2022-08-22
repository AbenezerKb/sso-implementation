Feature: Create user
    Background: I am logged in
        Given I am logged in with the following credentials
            | email           | password | role        |
            | test2@gmail.com | 1234567  | create_user |
    Scenario: Successfully create user
        When I fill the form with the following details
            | first_name | middle_name | last_name | phone      | email            | role   |
            | testuser1  | testuser1   | testuser1 | 0925252595 | test11@gmail.com | 123456 |
        And I submit the create user form
        Then The user is created

    Scenario Outline: Failed user creation
        When I fill the form with the following details
            | first_name   | middle_name   | last_name   | phone   | email   | role   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> | <role> |
        And I submit the create user form
        Then the creating process should fail with "<message>"

        Examples:
            | first_name | middle_name | last_name | phone      | email           | role   | message                 |
            |            | testuser1   | testuser1 | 0925252525 | test1@gmail.com | 123456 | first name is required  |
            | testuser1  |             | testuser1 | 0925252525 | test1@gmail.com | 123456 | middle name is required |
            | testuser1  | testuser1   |           | 0925252525 | test1@gmail.com | 123456 | last name is required   |
            | testuser1  | testuser1   | testuser1 |            | test1@gmail.com | 123456 | phone is required       |
            | testuser1  | testuser1   | testuser1 | 0925252525 | test1gmail.com  | 123456 | email is not valid      |
            | testuser1  | testuser1   | testuser1 | 33333333   | test1@gmail.com | 123456 | invalid phone number    |
            | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com |        | role is required        |