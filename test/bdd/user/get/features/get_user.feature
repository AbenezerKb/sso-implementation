Feature: Get User

    As an admin
    I want to get the user's details
    So that I can know the details of the particular user

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role          |
            | admin@gmail.com | iAmAdmin | get_user |

        And there is user with the following details:
            | first_name | middle_name | last_name | phone        | email           | password |
            | nati       | nati        | nati      | 251923456789 | normal@gmail.com | 123456   |
    @success
    Scenario: Successful Get user
        Given I have users id
        When I Get the user
        Then I should successfully get the user
    @failure
    Scenario Outline: user not found
        Given I have user with id "<id>"
        When I Get the user
        Then Then I should get error with message "<message>"

        Examples:
            | id                         | message        |
            | 3kjf0-kljkla0-afl30-afl-dk | user not found |
            

