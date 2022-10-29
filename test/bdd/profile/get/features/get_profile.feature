Feature: Get Profile

    As an user
    I want to get my profile details
    So that I can the see and modify my profile details

    Scenario: Successful Get Profile
        Given I am logged in user with the following details
            | first_name | middle_name | last_name | phone        | email            | password | gender | role      |
            | nati       | nati        | nati      | 251923456789 | normal@gmail.com | 123456   | male   | not-admin |
        When I request to get my profile
        Then I should successfully get my profile
