Feature: Delete User
    Background: setup test seed
        Given I am logged in with the following credentials
            | email           | password | role        |
            | test2@gmail.com | 1234567  | delete_user |
        And I have a registered users
            | first_name | middle_name | last_name | phone      | email            | role   |
            | testuser1  | testuser1   | testuser1 | 0925252595 | test11@gmail.com | 123456 |
            | testuser2  | testuser2   | testuser2 | 0925252596 | test12@gmail.com | 123456 |
            | testuser3  | testuser3   | testuser3 | 0925252597 | test13@gmail.com | 123456 |
    @success
    Scenario Outline: I successfully Delete the user
        When I request to delete the user
        Then the user should be deleted
    @failer
    Scenario Outline: User Should not be deleted with in invalid ID
        When I request to delete the users with in "<invalid_ID>"
        Then The system user should get an error message "<error_message>"
        Examples:
            | invalid_ID                           | error_message      |
            | 000000000                            | invalid user input |
            | a8aa9217-83ae-4f33-bce4-6ba81cedf13e | user not found     |