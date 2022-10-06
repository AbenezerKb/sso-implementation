Feature: update User Status

    As a admin
    I want to update user's status
    So that I can activate or deactivate user as appropriate

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role        |
            | admin@gmail.com | iAmAdmin | update_user_status |
    @success
    Scenario Outline: Successful User Status Update
        Given there is user with the following details:
            | first_name   | middle_name   | last_name   | phone   | email   | password   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> | <password> |
        When I update the user's status to "<status>"
        Then the user status should update to "<status>"

        Examples:
            | first_name | middle_name | last_name | phone        | email           | password | status   |
            | testuser1  | testuser1   | testuser1 | 251925252525 | test1@gmail.com | 123456   | INACTIVE |
            | testuser1  | testuser1   | testuser1 | 251925252525 | test1@gmail.com | 123456   | ACTIVE   |

    @failure
    Scenario Outline: Failed User Status Update
        Given there is user with id "<id>"
        When I update the user's status to "<status>"
        Then Then I should get user not found error with message "<message>"

        Examples:
            | id                         | message        |
            | 3kjf0-kljkla0-afl30-afl-dk | user not found |
    @failure
    Scenario Outline: Invalid Status
        Given there is user with the following details:
            | first_name   | middle_name   | last_name   | phone   | email   | password   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> | <password> |
        When I update the user's status to "<status>"
        Then Then I should get error with message "<message>"

        Examples:
            | first_name | middle_name | last_name | phone        | email           | password | status    | message               |
            | testuser1  | testuser1   | testuser1 | 251925252525 | test1@gmail.com | 123456   | INACTIVED | must be a valid value |
            | testuser1  | testuser1   | testuser1 | 251925252525 | test1@gmail.com | 123456   |           | status is required    |
