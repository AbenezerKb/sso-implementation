Feature: Delete User
    Background: setup test seed
        Given I have a registered account on the system
            | first_name | middle_name | last_name | phone      |
            | testuser1  | testuser1   | testuser1 | 0925252595 |
        And  I am logged in with the following credentials
            | email           | password | role        |
            | test2@gmail.com | 1234567  | delete_user |
    @success
    Scenario Outline: I successfully Delete my account
        When I want to delete my account
        Then My account should be deleted